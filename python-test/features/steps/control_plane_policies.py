from hamcrest import *
import requests
from behave import given, when, then, step
from utils import random_string, filter_list_by_parameter_start_with, safe_load_json, remove_empty_from_json
from local_agent import get_orb_agent_logs
from test_config import TestConfig
import time
from datetime import datetime
from control_plane_datasets import create_new_dataset, list_datasets
from random import choice, choices, sample

policy_name_prefix = "test_policy_name_"
base_orb_url = TestConfig.configs().get('base_orb_url')


@when("a new policy is created using: {kwargs}")
def create_new_policy(context, kwargs):

    acceptable_keys = ['policy_name', 'handler_label', 'handler', 'description', 'tap', 'input_type',
                       'host_specification', 'bpf_filter_expression', 'pcap_source', 'only_qname_suffix',
                       'only_rcode', 'backend_type']
    policy_name = policy_name_prefix + random_string(10)

    kwargs_dict = {'policy_name': policy_name, 'handler': None, 'description': None, 'tap': "default_pcap",
                   'input_type': "pcap", 'host_specification': None, 'bpf_filter_expression': None,
                   'pcap_source': None, 'only_qname_suffix': None, 'only_rcode': None, 'backend_type': "pktvisor"}

    for i in kwargs.split(", "):
        assert_that(i, matches_regexp("^.+=.+$"), f"Unexpected format for param {i}")
        item = i.split("=")
        kwargs_dict[item[0]] = item[1]

    assert_that(all(key in acceptable_keys for key, value in kwargs_dict.items()), equal_to(True),
                f"Unexpected parameters for policy. Options are {acceptable_keys}")

    if kwargs_dict["only_qname_suffix"] is not None:
        kwargs_dict["only_qname_suffix"] = kwargs_dict["only_qname_suffix"].replace("[", "")
        kwargs_dict["only_qname_suffix"] = kwargs_dict["only_qname_suffix"].replace("]", "")
        kwargs_dict["only_qname_suffix"] = kwargs_dict["only_qname_suffix"].split("/ ")

    handle_label = f"default_{kwargs_dict['handler']}_{random_string(3)}"

    context.policy = create_policy(context.token, kwargs_dict["policy_name"], handle_label,
                                   kwargs_dict["handler"], kwargs_dict["description"], kwargs_dict["tap"],
                                   kwargs_dict["input_type"], kwargs_dict["host_specification"],
                                   kwargs_dict["bpf_filter_expression"], kwargs_dict["pcap_source"],
                                   kwargs_dict["only_qname_suffix"], kwargs_dict["only_rcode"],
                                   kwargs_dict["backend_type"])

    assert_that(context.policy['name'], equal_to(policy_name))
    if 'policies_created' in context:
        context.policies_created[context.policy['id']] = context.policy['name']
    else:
        context.policies_created = dict()
        context.policies_created[context.policy['id']] = context.policy['name']


@then("referred policy {condition} be listed on the orb policies list")
def check_policies(context, condition='must'):
    policy_id = context.policy['id']
    all_existing_policies = list_policies(context.token)
    is_policy_listed = bool()
    for policy in all_existing_policies:
        if policy_id in policy.values():
            is_policy_listed = True
            break
        is_policy_listed = False
    if condition == 'must':
        assert_that(is_policy_listed, equal_to(True), f"Policy {policy_id} not listed on policies list")
        get_policy(context.token, policy_id)
    elif condition == 'must not':
        assert_that(is_policy_listed, equal_to(False), f"Policy {policy_id} exists in the policies list")
        policy = get_policy(context.token, policy_id, 404)
        assert_that(policy['error'], equal_to('non-existent entity'),
                    "Unexpected response for get policy request")


@step('one of applied policies is removed')
def remove_policy_applied(context):
    context.considered_timestamp = datetime.now().timestamp()
    policy_removed = choice(context.list_agent_policies_id)
    context.policy = get_policy(context.token, policy_removed)
    delete_policy(context.token, context.policy["id"])
    context.list_agent_policies_id.remove(context.policy["id"])
    context.policies_created.pop(context.policy["id"])
    existing_datasets = list_datasets(context.token)
    context.id_of_datasets_related_to_removed_policy = list_datasets_for_a_policy(policy_removed, existing_datasets)


@step('container logs should inform that removed policy was stopped and removed within {time_to_wait} seconds')
def check_test(context, time_to_wait):
    time_waiting = 0
    sleep_time = 0.5
    timeout = int(time_to_wait)
    found = {'stop': False, 'remove': False}
    stop_log_info = f"policy [{context.policy['name']}]: stopping"
    remove_log_info = f"DELETE /api/v1/policies/{context.policy['name']} 200"

    while time_waiting < timeout and (found['stop'] is False or found['remove'] is False):
        logs = get_orb_agent_logs(context.container_id)
        for log_line in logs:
            log_line = safe_load_json(log_line)
            if found['stop'] is False:
                found['stop'] = is_expected_log_info_in_log_line(log_line, stop_log_info, context.considered_timestamp)

            if found['remove'] is False:
                found['remove'] = is_expected_log_info_in_log_line(log_line, remove_log_info,
                                                                   context.considered_timestamp)
            if found['stop'] is True and found['remove'] is True:
                break
        time.sleep(sleep_time)
        time_waiting += sleep_time


@then('cleanup policies')
def clean_policies(context):
    """
    Remove all policies starting with 'policy_name_prefix' from the orb

    :param context: Behave class that contains contextual information during the running of tests.
    """
    token = context.token
    policies_list = list_policies(token)
    policies_filtered_list = filter_list_by_parameter_start_with(policies_list, 'name', policy_name_prefix)
    delete_policies(token, policies_filtered_list)


@given("that a policy using: {kwargs} already exists")
def new_policy(context, kwargs):
    create_new_policy(context, kwargs)
    check_policies(context)


@step('the container logs that were output after {condition} does not contain the message "{text_to_match}" referred '
      'to deleted policy anymore')
def check_agent_logs_for_deleted_policies_considering_timestamp(context, condition, text_to_match):
    policies_have_expected_message = \
        check_agent_log_for_policies(text_to_match, 0, context.container_id, list(context.policy['id']),
                                     context.considered_timestamp)
    assert_that(len(policies_have_expected_message), equal_to(0),
                f"Message '{text_to_match}' for policy "
                f"'{context.policy['id']}: {context.policy['name']}'"
                f" present on logs even after removing policy!")


@step('the container logs that were output after {condition} contain the message "{'
      'text_to_match}" referred to each applied policy within {time_to_wait} seconds')
def check_agent_logs_for_policies_considering_timestamp(context, condition, text_to_match, time_to_wait):
    policies_data = list()
    policies_have_expected_message = \
        check_agent_log_for_policies(text_to_match, time_to_wait, context.container_id, context.list_agent_policies_id,
                                     context.considered_timestamp)
    if len(set(context.list_agent_policies_id).difference(policies_have_expected_message)) > 0:
        policies_without_message = set(context.list_agent_policies_id).difference(policies_have_expected_message)
        for policy in policies_without_message:
            policies_data.append(get_policy(context.token, policy))

    assert_that(policies_have_expected_message, equal_to(set(context.list_agent_policies_id)),
                f"Message '{text_to_match}' for policy "
                f"'{policies_data}'"
                f" was not found in the agent logs!")


@step('the container logs contain the message "{text_to_match}" referred to each policy within {'
      'time_to_wait} seconds')
def check_agent_logs_for_policies(context, text_to_match, time_to_wait):
    policies_have_expected_message = \
        check_agent_log_for_policies(text_to_match, time_to_wait, context.container_id, context.list_agent_policies_id)
    assert_that(policies_have_expected_message, equal_to(set(context.list_agent_policies_id)),
                f"Message '{text_to_match}' for policy "
                f"'{set(context.list_agent_policies_id).difference(policies_have_expected_message)}'"
                f" was not found in the agent logs!")


@step('{amount_of_policies} {type_of_policies} policies are applied to the group')
def apply_n_policies(context, amount_of_policies, type_of_policies):
    args_for_policies = return_policies_type(int(amount_of_policies), type_of_policies)
    for i in range(int(amount_of_policies)):
        create_new_policy(context, args_for_policies[i][1])
        check_policies(context)
        create_new_dataset(context)


@step('{amount_of_policies} {type_of_policies} policies are applied to the group by {amount_of_datasets} datasets each')
def apply_n_policies_x_times(context, amount_of_policies, type_of_policies, amount_of_datasets):
    for n in range(int(amount_of_policies)):
        args_for_policies = return_policies_type(int(amount_of_policies), type_of_policies)
        create_new_policy(context, args_for_policies[n][1])
        check_policies(context)
        for x in range(int(amount_of_datasets)):
            create_new_dataset(context)


def create_policy(token, policy_name, handler_label, handler, description=None, tap="default_pcap",
                  input_type="pcap", host_specification=None, bpf_filter_expression=None, pcap_source=None,
                  only_qname_suffix=None, only_rcode=None, backend_type="pktvisor"):
    """

    Creates a new policy in Orb control plane

    :param (str) token: used for API authentication
    :param (str) policy_name:  of the policy to be created
    :param (str) handler_label:  of the handler
    :param (str) handler: to be added
    :param (str) description: description of policy
    :param tap: named, host specific connection specifications for the raw input streams accessed by pktvisor
    :param input_type: this must reference a tap name, or application of the policy will fail
    :param (str) host_specification: Subnets (comma separated) which should be considered belonging to this host,
    in CIDR form. Used for ingress/egress determination, defaults to host attached to the network interface.
    :param (str) bpf_filter_expression: these decide exactly which data to summarize and expose for collection.
                                        Tcpdump compatible filter expression for limiting the traffic examined
                                        (with BPF). See https://www.tcpdump.org/manpages/tcpdump.1.html.
    :param (str) pcap_source: Packet capture engine to use. Defaults to best for platform.
                                                            Options: af_packet (linux only) or libpcap.
    :param (str) only_qname_suffix: Filter out any queries whose QName does not end in a suffix on the list
    :param (int) only_rcode: Filter out any queries which are not the given RCODE. Options:
                                                                                    "NOERROR": 0,
                                                                                    "NXDOMAIN": 3,
                                                                                    "REFUSED": 5,
                                                                                    "SERVFAIL": 2
    :param backend_type: Agent backend this policy is for. Cannot change once created. Default: pktvisor
    :return: (dict) a dictionary containing the created policy data
    """
    if only_rcode is not None: only_rcode = int(only_rcode)
    assert_that(pcap_source, any_of(equal_to(None), equal_to("af_packet"), equal_to("libpcap")),
                "Unexpected type of pcap_source")
    assert_that(only_rcode, any_of(equal_to(None), equal_to(0), equal_to(2), equal_to(3), equal_to(5)),
                "Unexpected type of only_rcode")

    json_request = {"name": policy_name,
                    "description": description,
                    "backend": backend_type,
                    "policy": {
                        "kind": "collection",
                        "input": {
                            "tap": tap,
                            "input_type": input_type,
                            "config": {
                                "host_spec": host_specification,
                                "pcap_source": pcap_source},
                            "filter": {"bpf": bpf_filter_expression}},
                        "handlers": {
                            "modules": {
                                handler_label: {
                                    "type": handler,
                                    "filter": {
                                        "only_qname_suffix": only_qname_suffix,
                                        "only_rcode": only_rcode
                                    }
                                }
                            }
                        }
                    }
                    }
    json_request = remove_empty_from_json(json_request.copy())
    headers_request = {'Content-type': 'application/json', 'Accept': '*/*', 'Authorization': token}

    response = requests.post(base_orb_url + '/api/v1/policies/agent', json=json_request, headers=headers_request)
    assert_that(response.status_code, equal_to(201),
                'Request to create policy failed with status=' + str(response.status_code))

    return response.json()


def get_policy(token, policy_id, expected_status_code=200):
    """
    Gets a policy from Orb control plane

    :param (str) token: used for API authentication
    :param (str) policy_id: that identifies policy to be fetched
    :param (int) expected_status_code: expected request's status code. Default:200.
    :returns: (dict) the fetched policy
    """

    get_policy_response = requests.get(base_orb_url + '/api/v1/policies/agent/' + policy_id,
                                       headers={'Authorization': token})

    assert_that(get_policy_response.status_code, equal_to(expected_status_code),
                'Request to get policy id=' + policy_id + ' failed with status=' + str(get_policy_response.status_code))

    return get_policy_response.json()


def list_policies(token, limit=100):
    """
    Lists all policies from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :param (int) limit: Size of the subset to retrieve. (max 100). Default = 100
    :returns: (list) a list of policies
    """
    response = requests.get(base_orb_url + '/api/v1/policies/agent', headers={'Authorization': token},
                            params={'limit': limit})

    assert_that(response.status_code, equal_to(200),
                'Request to list policies failed with status=' + str(response.status_code))

    policies_as_json = response.json()
    return policies_as_json['data']


def delete_policies(token, list_of_policies):
    """
    Deletes from Orb control plane the policies specified on the given list

    :param (str) token: used for API authentication
    :param (list) list_of_policies: that will be deleted
    """

    for policy in list_of_policies:
        delete_policy(token, policy['id'])


def delete_policy(token, policy_id):
    """
    Deletes a policy from Orb control plane

    :param (str) token: used for API authentication
    :param (str) policy_id: that identifies the policy to be deleted
    """

    response = requests.delete(base_orb_url + '/api/v1/policies/agent/' + policy_id,
                               headers={'Authorization': token})

    assert_that(response.status_code, equal_to(204), 'Request to delete policy id='
                + policy_id + ' failed with status=' + str(response.status_code))


def check_logs_contain_message_for_policies(logs, expected_message, list_agent_policies_id, considered_timestamp):
    """
    Checks agent container logs for expected message for all applied policies and the log analysis loop is interrupted
    as soon as a log is found with the expected message for each applied policy.

    :param (list) logs: list of log lines
    :param (str) expected_message: message that we expect to find in the logs
    :param (list) list_agent_policies_id: list with all policy id applied to the agent
    :param (float) considered_timestamp: timestamp from which the log will be considered
    :returns: (set) set containing the ids of the policies for which the expected logs exist



    """
    policies_have_expected_message = set()
    for log_line in logs:
        log_line = safe_load_json(log_line)
        if is_expected_msg_in_log_line(log_line, expected_message, list_agent_policies_id,
                                       considered_timestamp) is True:
            policies_have_expected_message.add(log_line['policy_id'])
            if set(list_agent_policies_id) == set(policies_have_expected_message):
                return policies_have_expected_message
    return policies_have_expected_message


def check_agent_log_for_policies(expected_message, time_to_wait, container_id, list_agent_policies_id,
                                 considered_timestamp=datetime.now().timestamp()):
    """
    Checks agent container logs for expected message for each applied policy over a period of time

    :param (str) expected_message: message that we expect to find in the logs
    :param time_to_wait: seconds to wait for the log
    :param (str) container_id: agent container id
    :param (list) list_agent_policies_id: list with all policy id applied to the agent
    :param (float) considered_timestamp: timestamp from which the log will be considered.
                                                                Default: timestamp at which behave execution is started
    """
    time_waiting = 0
    sleep_time = 0.5
    timeout = int(time_to_wait)
    policies_have_expected_message = set()
    while time_waiting < timeout:
        logs = get_orb_agent_logs(container_id)
        policies_have_expected_message = \
            check_logs_contain_message_for_policies(logs, expected_message, list_agent_policies_id,
                                                    considered_timestamp)
        if len(policies_have_expected_message) == len(list_agent_policies_id):
            break
        time.sleep(sleep_time)
        time_waiting += sleep_time
    return policies_have_expected_message


def is_expected_msg_in_log_line(log_line, expected_message, list_agent_policies_id, considered_timestamp):
    """
    Test if log line has expected message
    - not be None
    - have a 'msg' property that matches the expected_message string.
    - have a 'ts' property whose value is greater than considered_timestamp
    - have a property 'policy_id' that is also contained in the list_agent_policies_id list

    :param (dict) log_line: agent container log line
    :param (str) expected_message: message that we expect to find in the logs
    :param (list) list_agent_policies_id: list with all policy id applied to the agent
    :param (float) considered_timestamp: timestamp from which the log will be considered.
    :return: (bool) whether expected message was found in the logs for expected policies

    """
    if log_line is not None:
        if log_line['msg'] == expected_message and 'policy_id' in log_line.keys():
            if log_line['policy_id'] in list_agent_policies_id:
                if log_line['ts'] > considered_timestamp:
                    return True
    return False


def is_expected_log_info_in_log_line(log_line, expected_log_info, considered_timestamp):
    """
    Test if log line has expected log
    - not be None
    - have a 'log' property that contains the expected_log_info string.
    - have a 'ts' property whose value is greater than considered_timestamp

    :param (dict) log_line: agent container log line
    :param (str) expected_log_info: log info that we expect to find in the logs
    :param (float) considered_timestamp: timestamp from which the log will be considered.
    :return: (bool) whether expected log info was found in the logs

    """
    if log_line is not None and 'log' in log_line.keys() and log_line['ts'] > considered_timestamp:
        if expected_log_info in log_line['log']:
            return True
    return False


def list_datasets_for_a_policy(policy_id, datasets_list):
    """

    :param (str) policy_id: that identifies the policy
    :param (list) datasets_list: list of datasets that will be filtered by policy
    :return: (list) list of ids of datasets related to referred policy
    """
    id_of_related_datasets = list()
    for dataset in datasets_list:
        if dataset['agent_policy_id'] == policy_id:
            id_of_related_datasets.append(dataset['id'])
    return id_of_related_datasets


def return_policies_type(k, policies_type='mixed'):

    assert_that(policies_type, any_of(equal_to('mixed'), any_of('simple'), any_of('advanced')),
                "Unexpected value for policies type")

    advanced = {
        'advanced_dns_libpcap_0': "handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=0",
        'advanced_dns_libpcap_2': "handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=2",
        'advanced_dns_libpcap_3': "handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=3",
        'advanced_dns_libpcap_5': "handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=5",

        'advanced_net': "handler=net, description='policy_net', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap",

        'advanced_dhcp': "handler=dhcp, description='policy_dhcp', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap",
        }

    simple = {

        'simple_dns': "handler=dns",
        'simple_net': "handler=net",
        # 'simple_dhcp': "handler=dhcp",
    }

    mixed = dict()
    mixed.update(advanced)
    mixed.update(simple)

    master_dict = {'advanced': advanced, 'simple': simple, 'mixed': mixed}

    if k <= len(master_dict[policies_type]):
        return sample(list(master_dict[policies_type].items()), k=k)

    return choices(list(master_dict[policies_type].items()), k=k)
