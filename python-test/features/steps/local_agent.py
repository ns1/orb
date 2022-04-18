from utils import safe_load_json, random_string
from behave import then, step
from hamcrest import *
from test_config import TestConfig, LOCAL_AGENT_CONTAINER_NAME
import docker
import time
import subprocess
import shlex
import threading

configs = TestConfig.configs()
ignore_ssl_and_certificate_errors = configs.get('ignore_ssl_and_certificate_errors')


@step('the agent container is started on port {port}')
def run_local_agent_container(context, port):
    if port.isdigit():
        port = int(port)
    assert_that(port, any_of(equal_to('default'), instance_of(int)), "Unexpected value for port")
    orb_address = configs.get('orb_address')
    interface = configs.get('orb_agent_interface', 'mock')
    agent_docker_image = configs.get('agent_docker_image', 'ns1labs/orb-agent')
    image_tag = ':' + configs.get('agent_docker_tag', 'latest')
    agent_image = agent_docker_image + image_tag
    env_vars = {"ORB_CLOUD_ADDRESS": orb_address,
                "ORB_CLOUD_MQTT_ID": context.agent['id'],
                "ORB_CLOUD_MQTT_CHANNEL_ID": context.agent['channel_id'],
                "ORB_CLOUD_MQTT_KEY": context.agent_key,
                "PKTVISOR_PCAP_IFACE_DEFAULT": interface}
    if ignore_ssl_and_certificate_errors == 'true':
        env_vars["ORB_TLS_VERIFY"] = "false"
    if port == "default":
        port = str(10853)
    else:
        env_vars["ORB_BACKENDS_PKTVISOR_API_PORT"] = str(port)

    context.container_id = run_agent_container(agent_image, env_vars, LOCAL_AGENT_CONTAINER_NAME)
    if port not in context.containers_id.keys():
        context.containers_id[str(port)] = context.container_id


@step('the container logs that were output after {condition} contain the message "{text_to_match}" within'
      '{time_to_wait} seconds')
def check_agent_logs_considering_timestamp(context, condition, text_to_match, time_to_wait):
    event = threading.Event()
    time_waiting = 0
    wait_time = 0.5
    timeout = int(time_to_wait)
    text_found = False

    while not event.is_set() and time_waiting < timeout:
        logs = get_orb_agent_logs(context.container_id)
        text_found = check_logs_contain_message(logs, text_to_match, event, context.considered_timestamp)
        if text_found is True:
            break
        event.wait(wait_time)
        time_waiting += wait_time

    assert_that(text_found, is_(True), 'Message "' + text_to_match + '" was not found in the agent logs!')


@then('the container logs should contain the message "{text_to_match}" within {time_to_wait} seconds')
def check_agent_log(context, text_to_match, time_to_wait):
    event = threading.Event()
    time_waiting = 0
    wait_time = 0.5
    timeout = int(time_to_wait)
    text_found = False

    while not event.is_set() and time_waiting < timeout:
        logs = get_orb_agent_logs(context.container_id)
        text_found = check_logs_contain_message(logs, text_to_match, event)
        if text_found is True:
            break
        time.sleep(wait_time)
        time_waiting += wait_time

    assert_that(text_found, is_(True), 'Message "' + text_to_match + '" was not found in the agent logs!')


@then("container on port {port} is {status} after {seconds} seconds")
def check_container_on_port_status(context, port, status, seconds):
    if port.isdigit():
        port = int(port)
    assert_that(port, any_of(equal_to('default'), instance_of(int)), "Unexpected value for port")
    if port == "default":
        port = str(10853)
    time.sleep(int(seconds))
    check_container_status(context.containers_id[str(port)], status)


@step("last container created is {status} after {seconds} seconds")
def check_last_container_status(context, status, seconds):
    time.sleep(int(seconds))
    check_container_status(context.container_id, status)


@step("the agent container is started using the command provided by the UI on port {port}")
def run_container_using_ui_command(context, port):
    if port.isdigit():
        port = int(port)
    assert_that(port, any_of(equal_to('default'), instance_of(int)), "Unexpected value for port")
    context.container_id = run_local_agent_from_terminal(context.agent_provisioning_command,
                                                         ignore_ssl_and_certificate_errors, str(port))
    assert_that(context.container_id, is_not((none())))
    rename_container(context.container_id, LOCAL_AGENT_CONTAINER_NAME)
    if port == "default":
        port = str(10853)
    if port not in context.containers_id.keys():
        context.containers_id[str(port)] = context.container_id


def run_agent_container(container_image, env_vars, container_name):
    """
    Gets a specific agent from Orb control plane

    :param (str) container_image: that will be used for running the container
    :param (dict) env_vars: that will be passed to the container context
    :param (str) container_name: base of container name
    :returns: (str) the container ID
    """
    LOCAL_AGENT_CONTAINER_NAME = container_name + random_string(5)
    client = docker.from_env()
    container = client.containers.run(container_image, name=LOCAL_AGENT_CONTAINER_NAME, detach=True,
                                      network_mode='host', environment=env_vars)
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


def check_logs_contain_message(logs, expected_message, event, start_time=0):
    """
    Gets the logs from Orb agent container

    :param (list) logs: list of log lines
    :param (str) expected_message: message that we expect to find in the logs
    :param (int) start_time: time to be considered as initial time. Default: None
    :returns: (bool) whether expected message was found in the logs
    """

    for log_line in logs:
        log_line = safe_load_json(log_line)

        if log_line is not None and log_line['msg'] == expected_message and log_line['ts'] > start_time:
            event.set()
            return event.is_set()

    return event.is_set()


def run_local_agent_from_terminal(command, ignore_ssl_and_certificate_errors, pktvisor_port):
    """
    :param (str) command: docker command to provision an agent
    :param (bool) ignore_ssl_and_certificate_errors: True if orb address doesn't have a valid certificate.
    :param (str or int) pktvisor_port: Port on which pktvisor should run
    :return: agent container ID
    """
    command = command.replace("\\\n", " ")
    args = shlex.split(command)
    if ignore_ssl_and_certificate_errors == 'true':
        args.insert(-1, "-e")
        args.insert(-1, "ORB_TLS_VERIFY=false")
    if pktvisor_port != 'default':
        args.insert(-1, "-e")
        args.insert(-1, f"ORB_BACKENDS_PKTVISOR_API_PORT={pktvisor_port}")
    terminal_running = subprocess.Popen(
        args, stdout=subprocess.PIPE)
    subprocess_return = terminal_running.stdout.read().decode()
    container_id = subprocess_return.split()
    assert_that(container_id[0], is_not((none())))
    return container_id[0]


def rename_container(container_id, container_name):
    """

    :param container_id: agent container ID
    :param container_name: base of agent container name
    """
    container_name = container_name + random_string(5)
    rename_container_command = f"docker rename {container_id} {container_name}"
    rename_container_args = shlex.split(rename_container_command)
    subprocess.Popen(rename_container_args, stdout=subprocess.PIPE)


def check_container_status(container_id, status):
    """

    :param container_id: agent container ID
    :param status: status that we expect to find in the container
    """
    docker_client = docker.from_env()
    container = docker_client.containers.list(all=True, filters={'id': container_id})
    assert_that(container, has_length(1))
    container = container[0]
    assert_that(container.status, equal_to(status), f"Container {container_id} failed with status {container.status}")
