from utils import safe_load_json, random_string, threading_wait_until, return_port_by_availability
from behave import then, step
from hamcrest import *
from configs import TestConfig, LOCAL_AGENT_CONTAINER_NAME
import docker
import subprocess
import shlex
from retry import retry
import threading
import json
from datetime import datetime
import ciso8601
from metrics import expected_metrics_by_handlers_and_groups, wait_until_metrics_scraped
from logger import Logger

log = Logger().logger_instance()

configs = TestConfig.configs()
verify_ssl_bool = eval(configs.get('verify_ssl').title())


@step("metrics must be correctly generated for {handler_type} handler")
def check_metrics_by_handler(context, handler_type):
    expected_metrics = expected_metrics_by_handlers_and_groups(handler_type, context.metric_groups_enabled,
                                                               context.metric_groups_disabled)
    policy_name = context.policy['name']
    local_prometheus_endpoint = f"http://localhost:{context.port}/api/v1/policies/{policy_name}/metrics/prometheus"
    correct_metrics, metrics_dif, metrics_present = wait_until_metrics_scraped(local_prometheus_endpoint,
                                                                               expected_metrics, timeout=60)
    expected_metrics_not_present = expected_metrics.difference(metrics_present)
    if expected_metrics_not_present == set():
        expected_metrics_not_present = None
    else:
        expected_metrics_not_present = sorted(expected_metrics_not_present)
    extra_metrics_present = metrics_present.difference(expected_metrics)
    if extra_metrics_present == set():
        extra_metrics_present = None
    else:
        extra_metrics_present = sorted(extra_metrics_present)
    assert_that(correct_metrics, equal_to(True), f"Metrics are not the expected. "
                                                 f"Metrics expected that are not present: {expected_metrics_not_present}."
                                                 f"Extra metrics present: {extra_metrics_present}")


@step('the agent container is started on an {status_port} port')
def run_local_agent_container(context, status_port, **kwargs):
    use_orb_live_address_pattern = configs.get("use_orb_live_address_pattern")
    verify_ssl = configs.get('verify_ssl')
    env_vars = create_agent_env_vars_set(context.agent['id'], context.agent['channel_id'], context.agent_key,
                                         verify_ssl, use_orb_live_address_pattern)
    env_vars.update(kwargs)
    assert_that(status_port, any_of(equal_to("available"), equal_to("unavailable")), "Unexpected value for port")
    availability = {"available": True, "unavailable": False}
    agent_docker_image = configs.get('agent_docker_image', 'orbcommunity/orb-agent')
    image_tag = ':' + configs.get('agent_docker_tag', 'latest')
    agent_image = agent_docker_image + image_tag

    context.port = return_port_by_availability(context, availability[status_port])

    if context.port != 10583:
        env_vars["ORB_BACKENDS_PKTVISOR_API_PORT"] = str(context.port)

    context.container_id = run_agent_container(agent_image, env_vars, LOCAL_AGENT_CONTAINER_NAME + random_string(2) +
                                               context.agent['name'][-5:])
    if context.container_id not in context.containers_id.keys():
        context.containers_id[context.container_id] = str(context.port)
    if availability[status_port]:
        log = f"web server listening on localhost:{context.port}"
    else:
        log = f"unable to bind to localhost:{context.port}"
    agent_started, logs, log_line = get_logs_and_check(context.container_id, log, element_to_check="log")
    assert_that(agent_started, equal_to(True),
                f"Log {log} not found on agent logs. Agent Name: {context.agent['name']}."
                f"\n Logs:{logs}")


@step('the agent container is started on an {status_port} port and use {group} env vars')
def run_local_agents_with_extra_env_vars(context, status_port, group):
    group = group.upper()
    assert_that(group, any_of("PCAP", "NETFLOW", "SFLOW", "DNSTAP", "ALL", "OTEL:ENABLED"))
    vars_by_input = {
        "PCAP": {"PKTVISOR_PCAP_IFACE_DEFAULT": configs.get("orb_agent_interface", "auto")},
        "NETFLOW": {"PKTVISOR_NETFLOW": "true", "PKTVISOR_NETFLOW_PORT_DEFAULT": 9995},
        "SFLOW": {"PKTVISOR_SFLOW": "true", "PKTVISOR_SFLOW_PORT_DEFAULT": 9994},
        "DNSTAP": {"PKTVISOR_DNSTAP": "true", "PKTVISOR_DNSTAP_PORT_DEFAULT": 9990}
    }
    if group == "ALL":
        vars_by_input["ALL"] = dict()
        vars_by_input["ALL"].update(vars_by_input["PCAP"])
        vars_by_input["ALL"].update(vars_by_input["NETFLOW"])
        vars_by_input["ALL"].update(vars_by_input["SFLOW"])
        vars_by_input["ALL"].update(vars_by_input["DNSTAP"])

    run_local_agent_container(context, status_port, **vars_by_input[group])


@step('the container logs that were output after {condition} contain the message "{text_to_match}" within'
      '{time_to_wait} seconds')
def check_agent_logs_considering_timestamp(context, condition, text_to_match, time_to_wait):
    # todo improve the logic for timestamp
    if "reset" in condition:
        considered_timestamp = context.considered_timestamp_reset
    else:
        considered_timestamp = context.considered_timestamp
    text_found, logs, log_line = get_logs_and_check(context.container_id, text_to_match, considered_timestamp,
                                          timeout=time_to_wait)
    assert_that(text_found, is_(True), f"Message {text_to_match} was not found in the agent logs!. \n\n"
                                       f"Container logs: {json.dumps(logs, indent=4)}")


@step("the container logs should not contain any {type_of_message} message")
def check_errors_on_agent_logs(context, type_of_message):
    type_dict = {"error": '"level":"error"', "panic": 'panic:'}
    logs = get_orb_agent_logs(context.container_id)
    non_expected_logs = [log for log in logs if type_dict[type_of_message] in log]
    assert_that(len(non_expected_logs), equal_to(0), f"agents logs contain the following {type_of_message}: "
                                                     f"{non_expected_logs}. \n All logs: {logs}.")


@then('the container logs should contain the message "{text_to_match}" within {time_to_wait} seconds')
def check_agent_msg_in_logs(context, text_to_match, time_to_wait):
    text_found, logs, log_line = get_logs_and_check(context.container_id, text_to_match, timeout=time_to_wait)

    assert_that(text_found, is_(True), f"Message {text_to_match} was not found in the agent logs!. \n\n"
                                       f"Container logs: {json.dumps(logs, indent=4)}")


@then('the container logs should contain "{error_log}" as log within {time_to_wait} seconds')
def check_agent_log_in_logs(context, error_log, time_to_wait):
    error_log = error_log.replace(":port", f":{context.port}")
    text_found, logs, log_line = get_logs_and_check(context.container_id, error_log, element_to_check="log", timeout=time_to_wait)
    assert_that(text_found, is_(True), f"Log {error_log} was not found in the agent logs!. \n\n"
                                       f"Container logs: {json.dumps(logs, indent=4)}")


@step("{order} container created is {status} within {seconds} seconds")
def check_last_container_status(context, order, status, seconds):
    order_convert = {"first": 0, "last": -1, "second": 1}
    container = list(context.containers_id.keys())[order_convert[order]]
    container_status = check_container_status(container, status, timeout=seconds)
    assert_that(container_status, equal_to(status), f"Container {context.container_id} failed with status "
                                                    f"{container_status}")


@step("{order} container created is {status} after {seconds} seconds")
def check_last_container_status_after_time(context, order, status, seconds):
    event = threading.Event()
    event.wait(int(seconds))
    event.set()
    if event.is_set() is True:
        check_last_container_status(context, order, status, seconds)


@step("the agent container is started using the command provided by the UI on an {status_port} port")
def run_container_using_ui_command(context, status_port):
    assert_that(status_port, any_of(equal_to("available"), equal_to("unavailable")), "Unexpected value for port")
    availability = {"available": True, "unavailable": False}
    context.port = return_port_by_availability(context, availability[status_port])
    verify_ssl = configs.get("verify_ssl")
    context.container_id = run_local_agent_from_terminal(context.agent_provisioning_command,
                                                         verify_ssl, str(context.port))
    assert_that(context.container_id, is_not((none())), f"Agent container was not run")
    rename_container(context.container_id, LOCAL_AGENT_CONTAINER_NAME + context.agent['name'][-5:])
    if context.container_id not in context.containers_id.keys():
        context.containers_id[context.container_id] = str(context.port)


@step(
    "the agent container is started using the command provided by the UI without {parameter_to_remove} on an {"
    "status_port} port")
def run_container_using_ui_command_without_restart(context, parameter_to_remove, status_port):
    context.agent_provisioning_command = context.agent_provisioning_command.replace(f"{parameter_to_remove} ", "")
    run_container_using_ui_command(context, status_port)


@step("stop the orb-agent container")
def stop_orb_agent_container(context):
    for container_id in context.containers_id.keys():
        stop_container(container_id)


@step("remove the orb-agent container")
def remove_orb_agent_container(context):
    for container_id in context.containers_id.keys():
        remove_container(container_id)
    context.containers_id = {}


@step("forced remove the orb-agent container")
def remove_orb_agent_container(context):
    for container_id in context.containers_id.keys():
        remove_container(container_id, force_remove=True)
    context.containers_id = {}


@step("force remove of all agent containers whose names start with the test prefix")
def remove_all_orb_agent_test_containers(context):
    docker_client = docker.from_env()
    containers = docker_client.containers.list(all=True)
    for container in containers:
        test_container = container.name.startswith(LOCAL_AGENT_CONTAINER_NAME)
        if test_container is True:
            container.remove(force=True)


def create_agent_env_vars_set(agent_id, agent_channel_id, agent_mqtt_key, verify_ssl,
                              use_orb_live_address_pattern):
    """
    Create the set of environmental variables to be passed to the agent
    :param agent_id: id of the agent
    :param agent_channel_id: id of the agent channel
    :param agent_mqtt_key: mqtt key to connect the agent
    :param verify_ssl: ignore process to verify tls if false
    :param use_orb_live_address_pattern: if true, uses the shortcut orb_cloud_address.
                                              if false sets api and mqtt address.
    :return: set of environmental variables
    """
    orb_address = configs.get('orb_address')
    env_vars = {"ORB_CLOUD_MQTT_ID": agent_id,
                "ORB_CLOUD_MQTT_CHANNEL_ID": agent_channel_id,
                "ORB_CLOUD_MQTT_KEY": agent_mqtt_key}
    if use_orb_live_address_pattern == "true":
        if orb_address != "orb.live":
            env_vars["ORB_CLOUD_ADDRESS"] = orb_address
        else:
            # default value must be enough to set correct parameters.
            pass
    else:
        env_vars["ORB_CLOUD_API_ADDRESS"] = configs.get("orb_url")
        env_vars["ORB_CLOUD_MQTT_ADDRESS"] = configs.get('mqtt_url')

    if verify_ssl == 'false':
        env_vars["ORB_TLS_VERIFY"] = "false"
    return env_vars


def run_agent_container(container_image, env_vars, container_name, time_to_wait=5):
    """
    Gets a specific agent from Orb control plane

    :param (str) container_image: that will be used for running the container
    :param (dict) env_vars: that will be passed to the container context
    :param (str) container_name: base of container name
    :param (int) time_to_wait: seconds that threading must wait after run the agent
    :returns: (str) the container ID
    """
    client = docker.from_env()
    container = client.containers.run(container_image, name=container_name, detach=True,
                                      network_mode='host', environment=env_vars)
    threading.Event().wait(time_to_wait)
    return container.id


def get_orb_agent_logs(container_id):
    """
    Gets the logs from Orb agent container

    :param (str) container_id: specifying the orb agent container
    :returns: (list) of log lines
    """
    docker_client = docker.from_env()
    container = docker_client.containers.get(container_id)
    return container.logs().decode("utf-8").split("\n")


def check_logs_contain_entry(logs, element_to_check, expected_entry, start_time=0):
    """
    Check if the logs from Orb agent container contain a specific entry

    :param (list) logs: list of log lines
    :param (str) element_to_check: key to search in the logs
    :param (str) expected_entry: entry that we expect to find in the logs
    :param (int) start_time: time to be considered as the initial time. Default: 0
    :returns: (bool) whether the expected entry was found in the logs
    """
    for log_line in logs:
        log_line = safe_load_json(log_line)

        if log_line is not None and element_to_check in log_line.keys() and isinstance(log_line['ts'], (int, str)):
            log_timestamp = (
                log_line['ts']
                if isinstance(log_line['ts'], int)
                else datetime.timestamp(ciso8601.parse_datetime(log_line['ts']))
            )

            if expected_entry in log_line.get(element_to_check, '') and log_timestamp > start_time:
                return True, log_line

    return False, None


def run_local_agent_from_terminal(command, verify_ssl, pktvisor_port):
    """
    :param (str) command: docker command to provision an agent
    :param (bool) verify_ssl: False if orb address doesn't have a valid certificate.
    :param (str or int) pktvisor_port: Port on which pktvisor should run
    :return: agent container ID
    """
    command = command.replace("\\\n", " ")
    args = shlex.split(command)
    if verify_ssl == 'false':
        args.insert(-1, "-e")
        args.insert(-1, "ORB_TLS_VERIFY=false")
    if pktvisor_port != 'default':
        args.insert(-1, "-e")
        args.insert(-1, f"ORB_BACKENDS_PKTVISOR_API_PORT={pktvisor_port}")
    terminal_running = subprocess.Popen(
        args, stdout=subprocess.PIPE)
    subprocess_return = terminal_running.stdout.read().decode()
    container_id = subprocess_return.split()
    assert_that(container_id[0], is_not((none())), f"Failed to run the agent. Command used: {args}.")
    return container_id[0]


@retry(tries=5, delay=0.2)
def rename_container(container_id, container_name):
    """

    :param container_id: agent container ID
    :param container_name: base of agent container name
    """
    docker_client = docker.from_env()
    containers = docker_client.containers.list(all=True)
    is_container_up = any(container_id in container.id for container in containers)
    assert_that(is_container_up, equal_to(True), f"Container {container_id} not found")
    container_name = container_name + random_string(5)
    rename_container_command = f"docker rename {container_id} {container_name}"
    rename_container_args = shlex.split(rename_container_command)
    subprocess.Popen(rename_container_args, stdout=subprocess.PIPE)


@threading_wait_until
def check_container_status(container_id, status, event=None):
    """

    :param container_id: agent container ID
    :param status: status that we expect to find in the container
    :param event: threading.event
    :return status of the container
    """
    docker_client = docker.from_env()
    container = docker_client.containers.list(all=True, filters={'id': container_id})
    assert_that(container, has_length(1), f"unable to find container {container_id}.")
    container = container[0]
    if container.status == status:
        event.set()
    return container.status


@threading_wait_until
def get_logs_and_check(container_id, expected_message, start_time=0, element_to_check="msg", event=None):
    """

    :param container_id: agent container ID
    :param (str) expected_message: message that we expect to find in the logs
    :param (int) start_time: time to be considered as initial time. Default: None
    :param element_to_check: Part of the log to be validated. Options: "msg" and "log". Default: "msg".
    :param (obj) event: threading.event
    :return: (bool) if the expected message is found return True, if not, False
    """
    assert_that(element_to_check, any_of(equal_to("msg"), equal_to("log")), "Unexpected value for element to check.")
    logs = get_orb_agent_logs(container_id)
    message_found, log_line = check_logs_contain_entry(logs, element_to_check, expected_message, start_time)
    if message_found is True:
        event.set()
    return message_found, logs, log_line


def run_agent_config_file(agent_name, overwrite_default=False, only_file=False, config_file_path="/opt/orb",
                          time_to_wait=5):
    """
    Run an agent container using an agent config file

    :param agent_name: name of the orb agent
    :param only_file: is true copy only the file. If false, copy the directory
    :param (bool) overwrite_default: if True and only_file is False saves the agent as "agent.yaml". Else, save it with
    agent name
    :param config_file_path: path to agent config file
    :param time_to_wait: seconds that threading must wait after run the agent
    :return: agent container id
    """
    agent_docker_image = configs.get('agent_docker_image', 'orbcommunity/orb-agent')
    agent_image = f"{agent_docker_image}:{configs.get('agent_docker_tag', 'latest')}"
    local_orb_path = configs.get("local_orb_path")
    if only_file is True:
        if overwrite_default is True:
            volume = f"{local_orb_path}/{agent_name}.yaml:{config_file_path}/agent.yaml"
        else:
            volume = f"{local_orb_path}/{agent_name}.yaml:{config_file_path}/{agent_name}.yaml"
    else:
        volume = f"{local_orb_path}:{config_file_path}/"
    agent_command = f"{config_file_path}/{agent_name}.yaml"
    if overwrite_default is True:
        command = f"docker run -d -v {volume} --net=host {agent_image}"
    else:
        command = f"docker run -d -v {volume} --net=host {agent_image} run -c {agent_command}"
    log.debug(f"Run Agent Command: {command}")
    args = shlex.split(command)
    terminal_running = subprocess.Popen(args, stdout=subprocess.PIPE)
    subprocess_return = terminal_running.stdout.read().decode()
    container_id = subprocess_return.split()[0]
    rename_container(container_id, LOCAL_AGENT_CONTAINER_NAME + agent_name[-5:])
    threading.Event().wait(time_to_wait)
    return container_id


def stop_container(container_id):
    """

    :param container_id: agent container ID
    """
    docker_client = docker.from_env()
    container = docker_client.containers.get(container_id)
    container.stop()


def remove_container(container_id, force_remove=False):
    """

    :param container_id: agent container ID
    :param force_remove: if True, similar to docker rm -f. Default: False
    """
    docker_client = docker.from_env()
    container = docker_client.containers.get(container_id)
    container.remove(force=force_remove)
