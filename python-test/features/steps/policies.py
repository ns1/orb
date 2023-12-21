from hamcrest import *
from behave import given, then, step
from utils import (random_string, filter_list_by_parameter_start_with, safe_load_json, remove_empty_from_json, \
                   threading_wait_until, UtilsManager, create_tags_set, is_json, values_to_boolean,
                   return_api_post_response,
                   return_api_get_response, return_api_put_response, return_api_delete_response)
from local_agent import get_orb_agent_logs
from configs import TestConfig
from datetime import datetime
from datasets import create_new_dataset, list_datasets
from random import choice, choices, sample
from deepdiff import DeepDiff
import json
import ciso8601
import yaml
from logger import Logger

log = Logger().logger_instance()
policy_name_prefix = "test_policy_name_"
configs = TestConfig.configs()
orb_url = configs.get('orb_url')
verify_ssl_bool = eval(configs.get('verify_ssl').title())


@step("a new policy is requested to be created with the same name as an existent one and: {kwargs}")
def create_policy_with_conflict_name(context, kwargs):
    if kwargs.split(", ")[-1].split("=")[-1] == "flow":
        kwargs_dict = parse_flow_policy_params(kwargs)
    else:
        kwargs_dict = parse_policy_params(kwargs)
    if kwargs_dict["handler"] == "flow":
        policy_json = make_policy_flow_json(context.policy['name'], kwargs_dict['handle_label'], kwargs_dict['handler'],
                                            kwargs_dict['description'],
                                            kwargs_dict['tap'], kwargs_dict['input_type'], kwargs_dict['port'],
                                            kwargs_dict['bind'], kwargs_dict['flow_type'],
                                            kwargs_dict['sample_rate_scaling'], kwargs_dict['only_devices'],
                                            kwargs_dict['only_ips'], kwargs_dict['only_ports'],
                                            kwargs_dict['only_interfaces'], kwargs_dict['geoloc_notfound'],
                                            kwargs_dict['asn_notfound'], kwargs_dict['backend_type'])
    else:
        policy_json = make_policy_json(context.policy['name'], kwargs_dict['handle_label'],
                                       kwargs_dict["handler"], kwargs_dict["description"], kwargs_dict["tap"],
                                       kwargs_dict["input_type"], kwargs_dict["host_specification"],
                                       kwargs_dict["bpf_filter_expression"], kwargs_dict["pcap_source"],
                                       kwargs_dict["only_qname_suffix"], kwargs_dict["only_rcode"],
                                       kwargs_dict["exclude_noerror"], kwargs_dict["backend_type"])

    context.error_message = create_policy(context.token, policy_json, expected_status_code=409)


@step("a {handler} policy {input_type} with tap_selector matching {match_type} tag(s) of the tap from {condition}, "
      "{metric_groups_enabled} metric_groups enabled, {metric_groups_disabled} metric_groups disabled and settings: {"
      "settings} is applied to the group")
def apply_policy_using_tap_selector(context, handler, input_type, match_type, condition, metric_groups_enabled,
                                    metric_groups_disabled, settings):
    module_name = f"{handler}_{random_string(5)}"
    policy_name = policy_name_prefix + random_string(10)
    if condition == "0 agent" and match_type == "any":
        tags = create_tags_set("3", tag_prefix='testtaptag', string_mode='lower')
    elif condition == "0 agent" and match_type == "all":
        tags = list(context.tap_tags.values())[0]
        tags.update(create_tags_set("1", tag_prefix='testtaptag', string_mode='lower'))
    elif condition == "1 agent (1 tag matching)":
        chosen_key = choice(list(context.tap_tags.keys()))
        tags = context.tap_tags[chosen_key]
    elif condition == "1 agent (1 tag matching + 1 random tag)":
        tags = create_tags_set("1", tag_prefix='testtaptag', string_mode='lower')
        chosen_key = choice(list(context.tap_tags.keys()))
        tags.update(context.tap_tags[chosen_key])
    elif condition == "an agent":
        tags = list(context.tap_tags.values())[0]
    else:
        raise ValueError("Invalid selector condition")

    policy = Policy(policy_name, f"description: {condition}", 'pktvisor')
    policy.add_input(input_type, 'tap_selector', input_match=match_type, tags=tags)
    if handler.lower() == "pcap":
        policy.add_pcap_module(module_name)
    elif handler.lower() == "dns":
        policy.add_dns_module(module_name, settings)
    elif handler.lower() == "net":
        policy.add_net_module(module_name, settings)
    elif handler.lower() == "dhcp":
        policy.add_dhcp_module(module_name)
    elif handler.lower() == "bgp":
        policy.add_bgp_module(module_name)
    elif handler.lower() == "flow":
        policy.add_flow_module(module_name, settings)
    elif handler.lower() == "netprobe":
        policy.add_netprobe_module(module_name)
    else:
        raise ValueError("Invalid policy handler. It must be one of pcap, dns, net, dhcp, bpg or flow.")
    if metric_groups_enabled.lower() != "default" and metric_groups_enabled.lower() != "none":
        policy.enable_metric_groups(module_name, metric_groups_enabled.split(", "))
    if metric_groups_disabled.lower() != "default" and metric_groups_disabled.lower() != "none":
        policy.disable_metric_groups(module_name, metric_groups_disabled.split(", "))
    json_for_create_policy = remove_empty_from_json(policy.policy)
    context.policy = create_policy(context.token, json_for_create_policy)
    check_policies(context)
    create_new_dataset(context, 1, 'last', 1, 'sink')
    context.metric_groups_enabled = metric_groups_enabled
    context.metric_groups_disabled = metric_groups_disabled


@step("the policy application error details must show that {message}")
def check_policy_error_detail(context, message):
    error_message = context.agent['last_hb_data']['policy_state'][context.policy['id']]['error']
    assert_that(message, equal_to(error_message), f"Unexpected error message. Agent: {context.agent}")


@step("a new policy is created using: {kwargs}")
def create_new_policy(context, kwargs):
    if kwargs.split(", ")[-1].split("=")[-1] == "flow":
        kwargs_dict = parse_flow_policy_params(kwargs)
    elif kwargs.split(", ")[-1].split("=")[-1] == "netprobe":
        kwargs_dict = parse_netprobe_policy_params(kwargs)
    else:
        kwargs_dict = parse_policy_params(kwargs)
    if kwargs_dict["handler"] == "flow":
        policy_json = make_policy_flow_json(kwargs_dict['name'], kwargs_dict['handle_label'], kwargs_dict['handler'],
                                            kwargs_dict['description'],
                                            kwargs_dict['tap'], kwargs_dict['input_type'], kwargs_dict['port'],
                                            kwargs_dict['bind'], kwargs_dict['flow_type'],
                                            kwargs_dict['sample_rate_scaling'], kwargs_dict['only_devices'],
                                            kwargs_dict['only_ips'], kwargs_dict['only_ports'],
                                            kwargs_dict['only_interfaces'], kwargs_dict['geoloc_notfound'],
                                            kwargs_dict['asn_notfound'], kwargs_dict['backend_type'])
    elif kwargs_dict["handler"] == "netprobe":
        policy_json = make_policy_netprobe_json(kwargs_dict["name"], kwargs_dict['handle_label'],
                                                kwargs_dict["handler"], kwargs_dict["description"], kwargs_dict["tap"],
                                                kwargs_dict["input_type"], kwargs_dict["test_type"],
                                                kwargs_dict["interval_msec"], kwargs_dict["timeout_msec"],
                                                kwargs_dict["packets_per_test"], kwargs_dict["packets_interval_msec"],
                                                kwargs_dict["packet_payload_size"], kwargs_dict["targets"],
                                                kwargs_dict["backend_type"])
    else:
        policy_json = make_policy_json(kwargs_dict["name"], kwargs_dict['handle_label'],
                                       kwargs_dict["handler"], kwargs_dict["description"], kwargs_dict["tap"],
                                       kwargs_dict["input_type"], kwargs_dict["host_specification"],
                                       kwargs_dict["bpf_filter_expression"], kwargs_dict["pcap_source"],
                                       kwargs_dict["only_qname_suffix"], kwargs_dict["only_rcode"],
                                       kwargs_dict["exclude_noerror"], kwargs_dict["backend_type"])

    context.policy = create_policy(context.token, policy_json)

    assert_that(context.policy['name'], equal_to(kwargs_dict["name"]), f"Policy name failed: {context.policy}")
    if 'policies_created' in context:
        context.policies_created[context.policy['id']] = context.policy['name']
    else:
        context.policies_created = dict()
        context.policies_created[context.policy['id']] = context.policy['name']


@step("editing a policy using {kwargs}")
def policy_editing(context, kwargs):
    acceptable_keys = ['name', 'handler_label', 'handler', 'description', 'tap', 'input_type',
                       'host_specification', 'bpf_filter_expression', 'pcap_source', 'only_qname_suffix',
                       'only_rcode', 'exclude_noerror', 'backend_type']

    handler_label = list(context.policy["policy"]["handlers"]["modules"].keys())[0]

    edited_attributes = {
        'host_specification': return_policy_attribute(context.policy, 'host_specification'),
        'bpf_filter_expression': return_policy_attribute(context.policy, 'bpf_filter_expression'),
        'pcap_source': return_policy_attribute(context.policy, 'pcap_source'),
        'only_qname_suffix': return_policy_attribute(context.policy, 'only_qname_suffix'),
        'only_rcode': return_policy_attribute(context.policy, 'only_rcode'),
        'description': return_policy_attribute(context.policy, 'description'),
        "name": return_policy_attribute(context.policy, 'name'),
        "handler": return_policy_attribute(context.policy, 'handler'),
        "backend_type": return_policy_attribute(context.policy, 'backend'),
        "tap": return_policy_attribute(context.policy, 'tap'),
        "input_type": return_policy_attribute(context.policy, 'input_type'),
        "handler_label": return_policy_attribute(context.policy, 'handler_label'),
        "exclude_noerror": return_policy_attribute(context.policy, "exclude_noerror")}

    if "host_spec" in context.policy["policy"]["input"]["config"].keys():
        edited_attributes["host_specification"] = context.policy["policy"]["input"]["config"]["host_spec"]
    if "pcap_source" in context.policy["policy"]["input"]["config"].keys():
        edited_attributes["pcap_source"] = context.policy["policy"]["input"]["config"]["pcap_source"]
    if "bpf" in context.policy["policy"]["input"]["filter"].keys():
        edited_attributes["bpf_filter_expression"] = context.policy["policy"]["input"]["filter"]["bpf"]
    if "description" in context.policy.keys():
        edited_attributes["description"] = context.policy['description']
    if "only_qname_suffix" in context.policy["policy"]["handlers"]["modules"][handler_label]['filter'].keys():
        edited_attributes["only_qname_suffix"] = \
            context.policy["policy"]["handlers"]["modules"][handler_label]["filter"][
                "only_qname_suffix"]
    if "only_rcode" in context.policy["policy"]["handlers"]["modules"][handler_label]['filter'].keys():
        edited_attributes["only_rcode"] = context.policy["policy"]["handlers"]["modules"][handler_label]["filter"][
            "only_rcode"]
    if "exclude_noerror" in context.policy["policy"]["handlers"]["modules"][handler_label]['filter'].keys():
        edited_attributes["exclude_noerror"] = context.policy["policy"]["handlers"]["modules"][handler_label]["filter"][
            "exclude_noerror"]

    for i in kwargs.split(", "):
        assert_that(i, matches_regexp("^.+=.+$"), f"Unexpected format for param {i}")
        item = i.split("=")
        edited_attributes[item[0]] = item[1]
        if item[1].isdigit() is False and str(item[1]).lower() == "none":
            edited_attributes[item[0]] = None
        if item[0] == "handler":
            edited_attributes["handler_label"] = f"default_{edited_attributes['handler']}_{random_string(3)}"

    for attribute in acceptable_keys:
        if attribute not in edited_attributes.keys():
            edited_attributes[attribute] = None

    assert_that(all(key in acceptable_keys for key, value in edited_attributes.items()), equal_to(True),
                f"Unexpected parameters for policy. Options are {acceptable_keys}")

    if edited_attributes["only_qname_suffix"] is not None:
        edited_attributes["only_qname_suffix"] = edited_attributes["only_qname_suffix"].replace("[", "")
        edited_attributes["only_qname_suffix"] = edited_attributes["only_qname_suffix"].replace("]", "")
        edited_attributes["only_qname_suffix"] = edited_attributes["only_qname_suffix"].split("/ ")

    if edited_attributes["name"] == 'conflict':
        policies_list = list_policies(context.token)
        policies_filtered_list = filter_list_by_parameter_start_with(policies_list, 'name', policy_name_prefix)
        policies_name = list()
        for policy in policies_filtered_list:
            policies_name.append(policy['name'])
        policies_name.remove(context.policy['name'])
        name_to_use = choice(policies_name)
        edited_attributes["name"] = name_to_use
        expected_status_code = 409
    else:
        expected_status_code = 200
        if policy_name_prefix not in edited_attributes["name"]:
            context.random_part_policy_name = f"_{random_string(10)}"
            edited_attributes["name"] = policy_name_prefix + edited_attributes["name"] + context.random_part_policy_name

    policy_json = make_policy_json(edited_attributes["name"], edited_attributes["handler_label"],
                                   edited_attributes["handler"], edited_attributes["description"],
                                   edited_attributes["tap"],
                                   edited_attributes["input_type"], edited_attributes["host_specification"],
                                   edited_attributes["bpf_filter_expression"], edited_attributes["pcap_source"],
                                   edited_attributes["only_qname_suffix"], edited_attributes["only_rcode"],
                                   edited_attributes["exclude_noerror"], edited_attributes["backend_type"])
    context.considered_timestamp = datetime.now().timestamp()

    if expected_status_code == 200:

        context.policy = edit_policy(context.token, context.policy['id'], policy_json)

        assert_that(context.policy['name'], equal_to(edited_attributes["name"]),
                    f"Policy name failed: {context.policy}")
    else:
        context.error_message = edit_policy(context.token, context.policy['id'], policy_json,
                                            expected_status_code=expected_status_code)


@step("policy {attribute} must be {value}")
def check_policy_attribute(context, attribute, value):
    acceptable_attributes = ['name', 'handler_label', 'handler', 'description', 'tap', 'input_type',
                             'host_specification', 'bpf_filter_expression', 'pcap_source', 'only_qname_suffix',
                             'only_rcode', 'backend_type', 'version', 'exclude_noerror']
    if attribute in acceptable_attributes:
        if attribute == "name":
            value = policy_name_prefix + value + context.random_part_policy_name
        policy_value = return_policy_attribute(context.policy, attribute)
        assert_that(str(policy_value), equal_to(value), f"Unexpected value for policy {attribute}")
    else:
        raise Exception(f"Attribute {attribute} not found on policy")


@step("otel policy version must be {policy_version}")
def check_otel_policy_version(context, policy_version):
    assert_that(policy_version.isdigit(), equal_to(True), "Policy version must be a number")
    policy_version = int(policy_version)
    policy_id = context.policy["id"]
    assert_that(context.policy.get("version"), equal_to(policy_version), f"Unexpected policy version in policy "
                                                                         f"{policy_id}")


@step("update otel policy using endpoint(s)={target_endpoint} and collection_interval={seconds}")
def update_otel_policy(context, target_endpoint, seconds):
    target_endpoints = target_endpoint.split(",")
    context.updated_policy_data = context.policy.get("policy")
    context.updated_policy_data["receivers"]["httpcheck"]["collection_interval"] = seconds
    targets = list()
    for endpoint in target_endpoints:
        targets.append(
            {'endpoint': endpoint, 'method': 'GET'}
        )
    context.updated_policy_data["receivers"]["httpcheck"]["targets"] = targets
    policy_yaml = yaml.dump(context.updated_policy_data)
    policy_json = make_otel_policy_json(context.policy['name'], context.policy.get("description"), policy_yaml)
    context.considered_timestamp = datetime.now().timestamp()
    context.policy = edit_policy(context.token, context.policy['id'], policy_json)


@step("updated fields were correctly present in updated otel policy")
def check_otel_policy_updated(context):
    context.policy = get_policy(context.token, context.policy['id'])
    assert_that(context.updated_policy_data["receivers"]["httpcheck"]["collection_interval"],
                equal_to(context.policy.get("policy")["receivers"]["httpcheck"]["collection_interval"]),
                f"Collection interval not updated. {context.policy}")
    assert_that(context.updated_policy_data["receivers"]["httpcheck"]["targets"],
                equal_to(context.policy.get("policy")["receivers"]["httpcheck"]["targets"]),
                f"Targets not updated. {context.policy}")


@then("referred policy {condition} be listed on the orb policies list")
def check_policies(context, **condition):
    if len(condition) > 0:
        condition = condition["condition"]
    else:
        condition = "must"
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
    if 'removed_policies_ids' in context:
        context.removed_policies_ids.append(context.policy["id"])
    else:
        context.removed_policies_ids = list()
        context.removed_policies_ids.append(context.policy["id"])
    context.list_agent_policies_id.remove(context.policy["id"])
    context.policies_created.pop(context.policy["id"])
    existing_datasets = list_datasets(context.token)
    context.id_of_datasets_related_to_removed_policy = list_datasets_for_a_policy(policy_removed, existing_datasets)


@step('container logs should inform that removed policy was stopped and removed within {time_to_wait} seconds')
def check_test(context, time_to_wait):
    stop_log_info = f"policy [{context.policy['name']}]: stopping"
    remove_log_info = f"DELETE /api/v1/policies/{context.policy['name']} 200"
    policy_removed = policy_stopped_and_removed(context.container_id, stop_log_info, remove_log_info,
                                                context.considered_timestamp, timeout=time_to_wait)
    assert_that(policy_removed, equal_to(True), f"Policy {context.policy} failed to be unapplied. \n"
                                                f"Agent: {json.dumps(context.agent, indent=4)}")


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
    policies_have_expected_message, logs = \
        check_agent_log_for_policies(text_to_match, context.container_id, list(context.policy['id']),
                                     context.considered_timestamp)
    assert_that(len(policies_have_expected_message), equal_to(0),
                f"Message '{text_to_match}' for policy "
                f"'{context.policy['id']}: {context.policy['name']}'"
                f" present on logs even after removing policy! \n"
                f"Agent: {json.dumps(context.agent, indent=4)}. \n"
                f"Agent Logs: {logs}")


@step('the container logs that were output after {condition} contain the message "{'
      'text_to_match}" referred to each applied policy within {time_to_wait} seconds')
def check_agent_logs_for_policies_considering_timestamp(context, condition, text_to_match, time_to_wait):
    # todo improve the logic for timestamp
    if "reset" in condition:
        considered_timestamp = context.considered_timestamp_reset
    else:
        considered_timestamp = context.considered_timestamp
    policies_data = list()
    policies_have_expected_message, logs = \
        check_agent_log_for_policies(text_to_match, context.container_id, context.list_agent_policies_id,
                                     considered_timestamp, timeout=time_to_wait)
    if len(set(context.list_agent_policies_id).difference(policies_have_expected_message)) > 0:
        policies_without_message = set(context.list_agent_policies_id).difference(policies_have_expected_message)
        for policy in policies_without_message:
            policies_data.append(get_policy(context.token, policy))
    log.debug(f"Agent logs: {logs}")
    assert_that(policies_have_expected_message, equal_to(set(context.list_agent_policies_id)),
                f"Message '{text_to_match}' for policy "
                f"'{policies_data}'"
                f" was not found in the agent logs!"
                f"Agent: {json.dumps(context.agent, indent=4)}. \n"
                f"Agent Logs: {logs}")


@step('the container logs contain the message "{text_to_match}" referred to each policy within {'
      'time_to_wait} seconds')
def check_agent_logs_for_policies(context, text_to_match, time_to_wait):
    policies_have_expected_message, logs = \
        check_agent_log_for_policies(text_to_match, context.container_id, context.list_agent_policies_id,
                                     timeout=time_to_wait)
    assert_that(policies_have_expected_message, equal_to(set(context.list_agent_policies_id)),
                f"Message '{text_to_match}' for policy "
                f"'{set(context.list_agent_policies_id).difference(policies_have_expected_message)}'"
                f" was not found in the agent logs!. \n"
                f"Agent: {json.dumps(context.agent, indent=4)}. \n"
                f"Agent Logs: {logs}")


@step('{amount_of_policies} policies with otel backend and yaml format are applied to the group')
def apply_n_otel_policies(context, amount_of_policies):
    policies = return_policies_otel(int(amount_of_policies))
    for policy_data in policies:
        name = f"{policy_name_prefix}otel_{random_string(10)}"
        description = "policy with otel backend"
        policy_yaml = yaml.dump(policy_data)
        policy_json = make_otel_policy_json(name, description, policy_yaml)
        context.policy = create_policy(context.token, policy_json)
        check_policies(context)
        if 'policies_created' in context:
            context.policies_created[context.policy['id']] = context.policy['name']
        else:
            context.policies_created = dict()
            context.policies_created[context.policy['id']] = context.policy['name']
        create_new_dataset(context, 1, 'last', 1, 'sink')


@step('{amount_of_policies} {type_of_policies} policies are applied to the group')
def apply_n_policies(context, amount_of_policies, type_of_policies):
    args_for_policies = return_policies_type(int(amount_of_policies), type_of_policies)
    for i in range(int(amount_of_policies)):
        create_new_policy(context, args_for_policies[i][1])
        check_policies(context)
        create_new_dataset(context, 1, 'last', 1, 'sink')


@step('{amount_of_policies} {type_of_policies} policies {policies_input} are applied to the group')
def apply_n_policies(context, amount_of_policies, type_of_policies, policies_input):
    if "same input_type as created via config file" in policies_input:
        policies_input = list(context.tap.values())[0]['input_type']
    args_for_policies = return_policies_type(int(amount_of_policies), type_of_policies, policies_input)
    if "tap" in context:
        tap_name = list(context.tap.keys())[0]
        input_type = list(context.tap.values())[0]['input_type']
    else:
        context.tap_name = tap_name = f"default_tap_before_provision_{random_string(10)}"
        input_type = policies_input
    for i in range(int(amount_of_policies)):
        kwargs = f"{args_for_policies[i][1]}, tap={tap_name}, input_type={input_type}"
        create_new_policy(context, kwargs)
        check_policies(context)
        create_new_dataset(context, 1, 'last', 1, 'sink')


@step('{amount_of_policies} {type_of_policies} policies are applied to the group by {amount_of_datasets} datasets each')
def apply_n_policies_x_times(context, amount_of_policies, type_of_policies, amount_of_datasets):
    for n in range(int(amount_of_policies)):
        args_for_policies = return_policies_type(int(amount_of_policies), type_of_policies)
        create_new_policy(context, args_for_policies[n][1])
        check_policies(context)
        for x in range(int(amount_of_datasets)):
            create_new_dataset(context, 1, 'last', 1, 'sink')


@step("{amount_of_policies} duplicated policies is applied to the group")
def apply_duplicate_policy(context, amount_of_policies):
    for i in range(int(amount_of_policies)):
        context.policy = create_duplicated_policy(context.token, context.policy["id"],
                                                  policy_name_prefix + random_string(10))
        check_policies(context)
        create_new_dataset(context, 1, 'last', 1, 'sink')


@step("try to duplicate this policy {times} times without set new name")
def duplicate_policy_with_same_name(context, times):
    # note that the context.policy is NOT changed, because we need to duplicate always the same policy to make the test
    # correctly
    context.duplicate_policies = list()
    for i in range(int(times)):
        if i <= 2:
            duplicated_policy = create_duplicated_policy(context.token, context.policy['id'])
        else:
            duplicated_policy = create_duplicated_policy(context.token, context.policy['id'], expected_status_code=409)
        context.duplicate_policies.append(duplicated_policy)


@step("try to duplicate this policy {times} times with a random new name")
def duplicate_policy_with_new_name(context, times):
    # note that the context.policy is NOT changed, because we need to duplicate always the same policy to make the test
    # correctly

    context.duplicate_policies = list()
    for i in range(int(times)):
        policy_new_name = policy_name_prefix + random_string(10)
        duplicated_policy = create_duplicated_policy(context.token, context.policy['id'],
                                                     new_policy_name=policy_new_name)
        context.duplicate_policies.append(duplicated_policy)


@step("{amount_successfully_policies} policies must be successfully duplicated and {amount_error_policies}"
      "must return an error")
def check_duplicated_policies_status(context, amount_successfully_policies, amount_error_policies):
    successfully_duplicated = list()
    wrongly_duplicated = 0
    for policy in context.duplicate_policies:
        if "id" in policy.keys():
            get_policy(context.token, policy['id'])
            successfully_duplicated.append(policy['id'])
        elif "error" in policy.keys():
            wrongly_duplicated += 1
    assert_that(len(successfully_duplicated),
                equal_to(int(amount_successfully_policies)), f"Amount of policies successfully duplicated fails."
                                                             f"Policies duplicated: {successfully_duplicated}")
    assert_that(wrongly_duplicated, equal_to(int(amount_error_policies)), f"Amount of policies wrongly duplicated fails"
                                                                          f".")


def create_duplicated_policy(token, policy_id, new_policy_name=None, expected_status_code=201):
    """

    :param (str) token: used for API authentication
    :param (str) policy_id: id of policy that will be duplicated
    :param (str) new_policy_name: name for the new policy created
    :param (int) expected_status_code: status code that must return on response
    :return: (dict) new policy created
    """
    json_request = {"name": new_policy_name}
    json_request = remove_empty_from_json(json_request)
    post_url = f"{orb_url}/api/v1/policies/agent/{policy_id}/duplicate"
    status_code, response = return_api_post_response(post_url, request_body=json_request, token=token,
                                                     verify=verify_ssl_bool)
    assert_that(status_code, equal_to(expected_status_code),
                'Request to create duplicated policy failed with status=' + str(status_code) + ': '
                + str(response))
    if status_code == 201:
        compare_two_policies(token, policy_id, response['id'])
    return response


def compare_two_policies(token, id_policy_one, id_policy_two):
    """

    :param (str) token: used for API authentication
    :param (str) id_policy_one: id of first policy
    :param str() id_policy_two: id of second policy

    """
    policy_one = get_policy(token, id_policy_one)
    policy_two = get_policy(token, id_policy_two)
    diff = DeepDiff(policy_one, policy_two, exclude_paths={"root['name']", "root['id']", "root['ts_last_modified']",
                                                           "root['ts_created']"})
    assert_that(diff, equal_to({}), f"Policy duplicated is not equal the one that generate it. Policy 1: {policy_one}\n"
                                    f"Policy 2: {policy_two}")


def create_policy(token, json_request, expected_status_code=201):
    """

    Creates a new policy in Orb control plane

    :param (str) token: used for API authentication
    :param (dict) json_request: policy json
    :expected_status_code (int): code to be returned on response
    :return: response of policy creation

    """

    status_code, response = return_api_post_response(f"{orb_url}/api/v1/policies/agent", request_body=json_request,
                                                     token=token, verify=verify_ssl_bool)

    assert_that(status_code, equal_to(expected_status_code),
                'Request to create policy failed with status=' + str(status_code) + ': '
                + str(response))

    return response


def edit_policy(token, policy_id, json_request, expected_status_code=200):
    """
    Editing a policy on Orb control plane

    :param (str) token: used for API authentication
    :param (str) policy_id: that identifies the policy to be edited
    :param (dict) json_request: policy json
    :param (int) expected_status_code: status to be returned on response
    :return: response of policy editing
    """

    status_code, response = return_api_put_response(f"{orb_url}/api/v1/policies/agent/{policy_id}",
                                                    request_body=json_request, token=token, verify=verify_ssl_bool)

    assert_that(status_code, equal_to(expected_status_code),
                'Request to editing policy failed with status=' + str(status_code) + ': '
                + str(response))

    return response


def make_policy_json(name, handler_label, handler, description=None, tap="default_pcap",
                     input_type="pcap", host_specification=None, bpf_filter_expression=None, pcap_source=None,
                     only_qname_suffix=None, only_rcode=None, exclude_noerror=None, backend_type="pktvisor"):
    """

    Generate a policy json

    :param (str) name:  of the policy to be created
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
    :param exclude_noerror: Filter out any queries which are not error response
    :param backend_type: Agent backend this policy is for. Cannot change once created. Default: pktvisor
    :return: (dict) a dictionary containing the created policy data
    """
    if only_rcode is not None: only_rcode = int(only_rcode)
    assert_that(pcap_source, any_of(equal_to(None), equal_to("af_packet"), equal_to("libpcap")),
                "Unexpected type of pcap_source")
    assert_that(only_rcode, any_of(equal_to(None), equal_to(0), equal_to(2), equal_to(3), equal_to(5)),
                "Unexpected type of only_rcode")
    if exclude_noerror is not None:
        assert_that(exclude_noerror.lower(), any_of(equal_to("false"), equal_to("true")),
                    "Unexpected value for exclude no error filter")
        exclude_noerror = eval(exclude_noerror.title())
    assert_that(handler, any_of(equal_to("dns"), equal_to("dhcp"), equal_to("net")), "Unexpected handler for policy")
    assert_that(name, not_none(), "Unable to create policy without name")

    if only_qname_suffix is not None and isinstance(only_qname_suffix, str):
        only_qname_suffix = only_qname_suffix.split(",")

    json_request = {"name": name,
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
                                        "only_rcode": only_rcode,
                                        "exclude_noerror": exclude_noerror
                                    }
                                }
                            }
                        }
                    }
                    }
    json_request = remove_empty_from_json(json_request.copy())
    return json_request


def make_otel_policy_json(name, description, policy_data):
    json_payload = {"name": name,
                    "backend": "otel",
                    "policy_data": policy_data,
                    "description": description,
                    "format": "yaml"}
    json_request = remove_empty_from_json(json_payload.copy())
    return json_request


def make_policy_flow_json(name, handler_label, handler, description=None, tap="default_flow",
                          input_type="flow", port=None, bind=None, flow_type=None, sample_rate_scaling=None,
                          only_devices=None, only_ips=None, only_ports=None, only_interfaces=None, geoloc_notfound=None,
                          asn_notfound=None, backend_type="pktvisor"):
    """

    Generate a policy json

    :param (str) name:  of the policy to be created
    :param (str) handler_label:  of the handler
    :param (str) handler: to be added
    :param (str) description: description of policy
    :param tap: named, host specific connection specifications for the raw input streams accessed by pktvisor
    :param input_type: this must reference a tap name, or application of the policy will fail
    :param backend_type: Agent backend this policy is for. Cannot change once created. Default: pktvisor
    :return: (dict) a dictionary containing the created policy data
    """
    assert_that(handler, equal_to("flow"), "Unexpected handler for policy")
    assert_that(name, not_none(), "Unable to create policy without name")

    json_request = {"name": name,
                    "description": description,
                    "backend": backend_type,
                    "policy": {
                        "kind": "collection",
                        "input": {
                            "tap": tap,
                            "input_type": input_type,
                            "config": {"port": port,
                                       "bind": bind,
                                       "only_ports": only_ports,
                                       "flow_type": flow_type}},
                        "handlers": {
                            "modules": {
                                handler_label: {
                                    "type": handler,
                                    "filter": {"only_devices": only_devices,
                                               "only_ips": only_ips,
                                               "only_ports": only_ports,
                                               "only_interfaces": only_interfaces,
                                               "geoloc_notfound": geoloc_notfound,
                                               "asn_notfound": asn_notfound},
                                    "config": {
                                        "sample_rate_scaling": sample_rate_scaling}
                                }
                            }
                        }
                    }
                    }
    json_request = remove_empty_from_json(json_request.copy())
    return json_request


def make_policy_netprobe_json(name, handler_label, handler, description=None, tap="default_netprobe",
                              input_type="flow", test_type='ping', interval_msec=None, timeout_msec=None,
                              packets_per_test=None, packets_interval_msec=None, packet_payload_size=None, targets=None,
                              backend_type="pktvisor"):
    """

    Generate a policy json

    :param (str) name:  of the policy to be created
    :param (str) handler_label:  of the handler
    :param (str) handler: to be added
    :param (str) description: description of policy
    :param tap: named, host specific connection specifications for the raw input streams accessed by pktvisor
    :param input_type: this must reference a tap name, or application of the policy will fail
    :param backend_type: Agent backend this policy is for. Cannot change once created. Default: pktvisor
    :return: (dict) a dictionary containing the created policy data
    """
    assert_that(handler, equal_to("netprobe"), "Unexpected handler for policy")
    assert_that(name, not_none(), "Unable to create policy without name")

    # netprobe configs are on tap level
    json_request = {"name": name,
                    "description": description,
                    "backend": backend_type,
                    "policy": {
                        "kind": "collection",
                        "input": {
                            "tap": tap,
                            "input_type": input_type,
                            "config": {"test_type": test_type,
                                       "interval_msec": interval_msec,
                                       "timeout_msec": timeout_msec,
                                       "packets_per_test": packets_per_test,
                                       "packets_interval_msec": packets_interval_msec,
                                       "packet_payload_size": packet_payload_size,
                                       "targets": targets}},
                        "handlers": {
                            "modules": {
                                handler_label: {
                                    "type": handler,
                                    "config": {},
                                    "filter": {}
                                }
                            }
                        }
                    }
                    }
    json_request = remove_empty_from_json(json_request.copy())
    return json_request


def get_policy(token, policy_id, expected_status_code=200):
    """
    Gets a policy from Orb control plane

    :param (str) token: used for API authentication
    :param (str) policy_id: that identifies policy to be fetched
    :param (int) expected_status_code: expected request's status code. Default:200.
    :returns: (dict) the fetched policy
    """

    status_code, response = return_api_get_response(f"{orb_url}/api/v1/policies/agent/{policy_id}", token=token,
                                                    verify=verify_ssl_bool)

    assert_that(status_code, equal_to(expected_status_code),
                'Request to get policy id=' + policy_id + ' failed with status= ' + str(status_code)
                + " response= " + str(response))

    return response


def list_policies(token, limit=100, offset=0):
    """
    Lists all policies from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :param (int) limit: Size of the subset to retrieve. (max 100). Default = 100
    :param (int) offset: Number of items to skip during retrieval. Default = 0.
    :returns: (list) a list of policies
    """

    all_policies, total, offset = list_up_to_limit_policies(token, limit, offset)

    new_offset = limit + offset

    while new_offset < total:
        policies_from_offset, total, offset = list_up_to_limit_policies(token, limit, new_offset)
        all_policies = all_policies + policies_from_offset
        new_offset = limit + offset

    return all_policies


def list_up_to_limit_policies(token, limit=100, offset=0):
    """
    Lists up to 100 policies from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :param (int) limit: Size of the subset to retrieve. (max 100). Default = 100
    :param (int) offset: Number of items to skip during retrieval. Default = 0.
    :returns: (list) a list of policies, (int) total policies on orb, (int) offset
    """

    status_code, response = return_api_get_response(f"{orb_url}/api/v1/policies/agent", token=token,
                                                    params={'limit': limit, 'offset': offset}, verify=verify_ssl_bool)

    assert_that(status_code, equal_to(200),
                'Request to list policies failed with status=' + str(status_code) + ': '
                + str(response))
    assert_that(response, has_key('data'), f"Response does not contain 'data' key: {response}")
    assert_that(response, has_key('total'), f"Response does not contain 'total' key: {response}")
    assert_that(response, has_key('offset'), f"Response does not contain 'offset' key: {response}")

    return response['data'], response['total'], response['offset']


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

    status_code, response = return_api_delete_response(f"{orb_url}/api/v1/policies/agent/{policy_id}",
                                                       token=token, verify=verify_ssl_bool)

    assert_that(status_code, equal_to(204), 'Request to delete policy id='
                + policy_id + ' failed with status=' + str(status_code))


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


@threading_wait_until
def check_agent_log_for_policies(expected_message, container_id, list_agent_policies_id,
                                 considered_timestamp=datetime.now().timestamp(), event=None):
    """
    Checks agent container logs for expected message for each applied policy over a period of time

    :param (str) expected_message: message that we expect to find in the logs
    :param (str) container_id: agent container id
    :param (list) list_agent_policies_id: list with all policy id applied to the agent
    :param (float) considered_timestamp: timestamp from which the log will be considered.
                                                                Default: timestamp at which behave execution is started
    :param (obj) event: threading.event
    """
    logs = get_orb_agent_logs(container_id)
    policies_have_expected_message = \
        check_logs_contain_message_for_policies(logs, expected_message, list_agent_policies_id,
                                                considered_timestamp)
    if len(policies_have_expected_message) == len(list_agent_policies_id):
        event.set()
        return policies_have_expected_message, logs

    return policies_have_expected_message, logs


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
        if expected_message in log_line['msg'] and 'policy_id' in log_line.keys():
            if log_line['policy_id'] in list_agent_policies_id:
                if isinstance(log_line['ts'], int) and log_line['ts'] > considered_timestamp:
                    return True
                elif isinstance(log_line['ts'], str) and datetime.timestamp(ciso8601.parse_datetime(log_line['ts'])) > \
                        considered_timestamp:
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
    if log_line is not None and 'log' in log_line.keys() and isinstance(log_line['ts'], int) and log_line['ts'] > \
            considered_timestamp:
        if expected_log_info in log_line['log']:
            return True
    elif log_line is not None and 'log' in log_line.keys() and isinstance(log_line['ts'], str) and \
            datetime.timestamp(ciso8601.parse_datetime(log_line['ts'])) > considered_timestamp:
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


def return_policies_otel(k):
    policies_otel = list()
    policies_otel.append({
        "receivers": {
            "httpcheck": {
                "targets": [
                    {
                        "endpoint": "http://orb.community",
                        "method": "GET"
                    },
                    {
                        "endpoint": "https://orb.live",
                        "method": "GET"
                    }
                ],
                "collection_interval": "60s"
            }
        },
        "extensions": None,
        "exporters": None,
        "service": {
            "pipelines": {
                "metrics": {
                    "exporters": None,
                    "receivers": [
                        "httpcheck"
                    ]
                }
            }
        }
    })
    policies_otel.append({
        "receivers": {
            "httpcheck": {
                "targets": [
                    {
                        "endpoint": "https://github.com/orb-community/orb",
                        "method": "GET"
                    }
                ],
                "collection_interval": "10s"
            }
        },
        "extensions": None,
        "exporters": None,
        "service": {
            "pipelines": {
                "metrics": {
                    "exporters": ["otlp"],
                    "receivers": [
                        "httpcheck"
                    ]
                }
            }
        }
    })
    if k <= len(policies_otel):
        return sample(policies_otel, k=k)

    return choices(policies_otel, k=k)


def return_policies_type(k, policies_type='mixed', input_type="pcap"):
    assert_that(policies_type, any_of(equal_to('mixed'), any_of('simple'), any_of('advanced')),
                "Unexpected value for policies type")

    if input_type == "flow":
        advanced = {
            "advanced_flow": "handler=flow, description='policy_flow', asn_notfound=true, sample_rate_scaling=true"
        }
        simple = {
            'simple_flow': "handler=flow"
        }
    elif input_type == "netprobe":
        advanced = {
            "advanced_netprobe_1": "handler=netprobe, test_type=ping, interval_msec=3000, timeout_msec=1000, packets_per_test=2, packets_interval_msec=30, packet_payload_size=56",
            "advanced_netprobe_2": "handler=netprobe, test_type=ping, packet_payload_size=56",
            "advanced_netprobe_3": "handler=netprobe, test_type=ping, interval_msec=900, timeout_msec=500, packets_per_test=5, packets_interval_msec=45"
        }
        simple = {
            'simple_netprobe': "handler=netprobe, test_type=ping"
        }
    else:
        advanced = {
            'advanced_dns_libpcap_0': "handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.orb.live/ .google.com], only_rcode=0",
            'advanced_dns_libpcap_2': "handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.orb.live/ .google.com], only_rcode=2",
            'advanced_dns_libpcap_3': "handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.orb.live/ .google.com], only_rcode=3",
            'advanced_dns_libpcap_5': "handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.orb.live/ .google.com], only_rcode=5",

            'advanced_net': "handler=net, description='policy_net', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap",
        }

        simple = {

            'simple_dns': "handler=dns",
            'simple_net': "handler=net"
        }

        if input_type != "dnstap":
            advanced[
                'advanced_dhcp'] = "handler=dhcp, description='policy_dhcp', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap"
            simple['simple_dhcp'] = "handler=dhcp"

    mixed = dict()
    mixed.update(advanced)
    mixed.update(simple)

    master_dict = {'advanced': advanced, 'simple': simple, 'mixed': mixed}

    if k <= len(master_dict[policies_type]):
        return sample(list(master_dict[policies_type].items()), k=k)

    return choices(list(master_dict[policies_type].items()), k=k)


def return_policy_attribute(policy, attribute):
    """

    :param (dict) policy: json of policy
    :param (str) attribute: policy attribute whose value is to be returned
    :return: (str, bool or None) value referring to policy attribute

    """

    handler_label = list(policy["policy"]["handlers"]["modules"].keys())[0]
    if attribute == "name":
        return policy['name']
    elif attribute == "handler_label":
        return handler_label
    elif attribute == "handler":
        return list(policy["policy"]["handlers"]["modules"].values())[0]["type"]
    elif attribute == "backend_type":
        return policy["backend"]
    elif attribute == "tap":
        return policy["policy"]["input"]["tap"]
    elif attribute == "input_type":
        return policy["policy"]["input"]["input_type"]
    elif attribute == "version" and "version" in policy.keys():
        return policy["version"]
    elif attribute == "description" and "description" in policy.keys():
        return policy['description']
    elif attribute == "host_specification" and "host_spec" in policy["policy"]["input"]["config"].keys():
        return policy["policy"]["input"]["config"]["host_spec"]
    elif attribute == "bpf_filter_expression" and "bpf" in policy["policy"]["input"]["filter"].keys():
        return policy["policy"]["input"]["filter"]["bpf"]
    elif attribute == "pcap_source" and "pcap_source" in policy["policy"]["input"]["config"].keys():
        return policy["policy"]["input"]["config"]["pcap_source"]
    elif attribute == "only_qname_suffix" and "only_qname_suffix" in \
            policy["policy"]["handlers"]["modules"][handler_label]["filter"].keys():
        return policy["policy"]["handlers"]["modules"][handler_label]["filter"]["only_qname_suffix"]
    elif attribute == "exclude_noerror" and "exclude_noerror" in \
            policy["policy"]["handlers"]["modules"][handler_label]["filter"].keys():
        return policy["policy"]["handlers"]["modules"][handler_label]["filter"]["exclude_noerror"]
    elif attribute == "only_rcode" and "only_rcode" in policy["policy"]["handlers"]["modules"][handler_label][
        "filter"].keys():
        return policy["policy"]["handlers"]["modules"][handler_label]["filter"]["only_rcode"]
    else:
        return None


@threading_wait_until
def policy_stopped_and_removed(container_id, stop_policy_info, remove_policy_info, start_considering_time, event=None):
    """

    :param (str) container_id: agent container id
    :param (str) stop_policy_info: log info that confirms that the policy was stopped
    :param (str) remove_policy_info: log info that confirms that the policy was removed
    :param (str) start_considering_time: timestamp after which logs must be validated
    :param (obj) event: threading.event
    :return: (bool) if the expected message is found return True, if not, False
    """
    found = {'stop': False, 'remove': False}
    logs = get_orb_agent_logs(container_id)
    for log_line in logs:
        log_line = safe_load_json(log_line)
        if found['stop'] is False:
            found['stop'] = is_expected_log_info_in_log_line(log_line, stop_policy_info, start_considering_time)

        if found['remove'] is False:
            found['remove'] = is_expected_log_info_in_log_line(log_line, remove_policy_info,
                                                               start_considering_time)
        if found['stop'] is True and found['remove'] is True:
            event.set()
            return event.is_set()
    return event.is_set()


def parse_policy_params(kwargs):
    acceptable_keys = ['name', 'handler_label', 'handler', 'description', 'tap', 'input_type',
                       'host_specification', 'bpf_filter_expression', 'pcap_source', 'only_qname_suffix',
                       'only_rcode', 'exclude_noerror', 'backend_type']

    name = policy_name_prefix + random_string(10)

    kwargs_dict = {'name': name, 'handler': None, 'description': None, 'tap': "default_pcap",
                   'input_type': "pcap", 'host_specification': None, 'bpf_filter_expression': None,
                   'pcap_source': None, 'only_qname_suffix': None, 'only_rcode': None, 'exclude_noerror': None,
                   'backend_type': "pktvisor"}

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

    if policy_name_prefix not in kwargs_dict["name"]:
        kwargs_dict["name"] + policy_name_prefix + kwargs_dict["name"]

    assert_that(kwargs_dict["handler"], any_of(equal_to("dns"), equal_to("dhcp"), equal_to("net")),
                "Unexpected handler for policy")
    kwargs_dict['handle_label'] = f"default_{kwargs_dict['handler']}_{random_string(3)}"

    return kwargs_dict


def parse_flow_policy_params(kwargs):
    name = policy_name_prefix + random_string(10)

    kwargs_dict = {'name': name, 'handler': None, 'description': None, 'tap': "default_flow",
                   'input_type': "flow", 'port': None, 'bind': None, 'flow_type': None, 'sample_rate_scaling': None,
                   'only_devices': None, 'only_ips': None, 'only_ports': None, 'only_interfaces': None,
                   'geoloc_notfound': None,
                   'asn_notfound': None, 'backend_type': "pktvisor"}

    for i in kwargs.split(", "):
        assert_that(i, matches_regexp("^.+=.+$"), f"Unexpected format for param {i}")
        item = i.split("=")
        kwargs_dict[item[0]] = item[1]

    if policy_name_prefix not in kwargs_dict["name"]:
        kwargs_dict["name"] + policy_name_prefix + kwargs_dict["name"]

    assert_that(kwargs_dict["handler"], equal_to("flow"), "Unexpected handler for policy")
    kwargs_dict['handle_label'] = f"default_{kwargs_dict['handler']}_{random_string(3)}"

    return kwargs_dict


def parse_netprobe_policy_params(kwargs):
    name = policy_name_prefix + random_string(10)

    kwargs_dict = {'name': name, 'handler': None, 'description': None, 'tap': "default_netprobe",
                   'input_type': "netprobe", 'test_type': 'ping', 'interval_msec': None, 'timeout_msec': None,
                   'packets_per_test': None, 'packets_interval_msec': None, 'packet_payload_size': None,
                   'targets': None, 'backend_type': "pktvisor"}

    for i in kwargs.split(", "):
        assert_that(i, matches_regexp("^.+=.+$"), f"Unexpected format for param {i}")
        item = i.split("=")
        kwargs_dict[item[0]] = item[1]

    if policy_name_prefix not in kwargs_dict["name"]:
        kwargs_dict["name"] + policy_name_prefix + kwargs_dict["name"]

    assert_that(kwargs_dict["handler"], equal_to("netprobe"), "Unexpected handler for policy")
    kwargs_dict['handle_label'] = f"default_{kwargs_dict['handler']}_{random_string(3)}"

    return kwargs_dict


class HandlerConfigs(UtilsManager):
    def __init__(self):
        self.handler_configs = dict()

    def add_configs(self, **kwargs):
        self.handler_configs = UtilsManager.add_configs(self, self.handler_configs, **kwargs)

        return self.handler_configs

    def remove_configs(self, *args):
        self.handler_configs = UtilsManager.remove_configs(self, self.handler_configs, *args)

        return self.handler_configs

    def json(self):
        return json.dumps(self.handler_configs)


class HandlerModules(HandlerConfigs):
    def __init__(self):
        self.handler_modules = dict()

    def __build_module(self, name, module_type, configs_list, filters_list, require_version=None):
        module = {
            name: {
                "type": module_type,
                "config": {
                },

                "filter": {
                },
                "metric_groups": {
                }
            }
        }
        if require_version is not None:
            module[name]["require_version"] = require_version

        module = UtilsManager.update_object_with_filters_and_configs(self, module, name, configs_list, filters_list)

        self.handler_modules.update(module)

    def __parse_module_settings(self, settings):
        if settings is None or settings == "default":
            settings_json = {}
        else:
            settings_is_json, settings_json = is_json(settings)
            assert_that(settings_is_json, is_(True), f"settings must be written in json format. Current settings: "
                                                     f"{settings}")
            settings_json = values_to_boolean(settings_json)
        return settings_json

    def add_dns_module(self, name, settings=None):

        settings_json = self.__parse_module_settings(settings)

        self.name = name
        self.public_suffix_list = {'public_suffix_list': settings_json.get('public_suffix_list', None)}
        self.only_rcode = {'only_rcode': settings_json.get("only_rcode", None)}
        self.exclude_noerror = {'exclude_noerror': settings_json.get("exclude_noerror", None)}
        self.only_dnssec_response = {'only_dnssec_response': settings_json.get("only_dnssec_response", None)}
        self.answer_count = {'answer_count': settings_json.get("answer_count", None)}
        self.only_qtype = {'only_qtype': settings_json.get("only_qtype", None)}
        self.only_qname_suffix = {'only_qname_suffix': settings_json.get("only_qname_suffix", None)}
        self.geoloc_notfound = {'geoloc_notfound': settings_json.get("geoloc_notfound", None)}
        self.asn_notfound = {'asn_notfound': settings_json.get("asn_notfound", None)}
        self.dnstap_msg_type = {'dnstap_msg_type': settings_json.get("dnstap_msg_type", None)}
        self.require_version = settings_json.get("require_version", None)

        dns_configs = [self.public_suffix_list]

        dns_filters = [self.only_rcode, self.exclude_noerror, self.only_dnssec_response, self.answer_count,
                       self.only_qtype, self.only_qname_suffix,
                       self.geoloc_notfound, self.asn_notfound, self.dnstap_msg_type]

        self.__build_module(self.name, "dns", dns_configs, dns_filters, self.require_version)
        return self.handler_modules

    def add_net_module(self, name, settings=None):

        settings_json = self.__parse_module_settings(settings)

        self.name = name
        self.geoloc_notfound = {'geoloc_notfound': settings_json.get('geoloc_notfound', None)}
        self.asn_notfound = {'asn_notfound': settings_json.get('asn_notfound', None)}
        self.only_geoloc_prefix = {'only_geoloc_prefix': settings_json.get('only_geoloc_prefix', None)}
        self.only_asn_number = {'only_asn_number': settings_json.get('only_asn_number', None)}
        self.require_version = settings_json.get("require_version", None)

        net_configs = []

        net_filters = [self.geoloc_notfound, self.asn_notfound, self.only_geoloc_prefix, self.only_asn_number]

        self.__build_module(self.name, "net", net_configs, net_filters, self.require_version)
        return self.handler_modules

    def add_dhcp_module(self, name):
        self.name = name

        dhcp_configs = []

        dhcp_filters = []

        self.__build_module(self.name, "dhcp", dhcp_configs, dhcp_filters)
        return self.handler_modules

    def add_bgp_module(self, name):
        self.name = name

        bgp_configs = []

        bgp_filters = []

        self.__build_module(self.name, "bgp", bgp_configs, bgp_filters)
        return self.handler_modules

    def add_pcap_module(self, name):
        self.name = name

        pcap_configs = []

        pcap_filters = []

        self.__build_module(self.name, "pcap", pcap_configs, pcap_filters)
        return self.handler_modules

    def add_flow_module(self, name, settings=None):

        settings_json = self.__parse_module_settings(settings)

        self.name = name
        self.sample_rate_scaling = {'sample_rate_scaling': settings_json.get("sample_rate_scaling", None)}
        self.recorded_stream = {'recorded_stream': settings_json.get("recorded_stream", None)}
        self.only_devices = {'only_devices': settings_json.get("only_devices", None)}
        self.only_ips = {'only_ips': settings_json.get("only_ips", None)}
        self.only_ports = {'only_ports': settings_json.get("only_ports", None)}
        self.only_interfaces = {'only_interfaces': settings_json.get("only_interfaces", None)}
        self.geoloc_notfound = {'geoloc_notfound': settings_json.get("geoloc_notfound", None)}
        self.asn_notfound = {'asn_notfound': settings_json.get("asn_notfound", None)}

        flow_configs = [self.sample_rate_scaling, self.recorded_stream]

        flow_filters = [self.only_devices, self.only_ips, self.only_ports, self.only_interfaces, self.geoloc_notfound,
                        self.asn_notfound]

        self.__build_module(self.name, "flow", flow_configs, flow_filters)
        return self.handler_modules

    def add_netprobe_module(self, name):
        self.name = name

        netprobe_configs = []

        netprobe_filters = []

        self.__build_module(self.name, "netprobe", netprobe_configs, netprobe_filters)
        return self.handler_modules

    def add_configs(self, name, **kwargs):
        self.handler_modules[name]["config"] = UtilsManager.add_configs(self, self.handler_modules[name]["config"],
                                                                        **kwargs)

        return self.handler_modules

    def add_filters(self, name, **kwargs):
        if "filter" not in self.handler_modules[name].keys():
            self.handler_modules[name].update({"filter": {}})

        self.handler_modules[name]["filter"] = UtilsManager.add_filters(self, self.handler_modules[name]["filter"],
                                                                        **kwargs)

        return self.handler_modules

    def enable_metric_groups(self, name, args):
        self.metric_groups = self.handler_modules[name]["metric_groups"]
        metrics_enable = list()
        if 'enable' not in self.metric_groups.keys():
            self.metric_groups.update({"enable": metrics_enable})

        for metric in args:
            metrics_enable.append(metric)
            if 'disable' in self.metric_groups.keys() and metric in self.metric_groups['disable']:
                self.metric_groups['disable'].remove(metric)

        self.metric_groups['enable'] = metrics_enable

        return self.handler_modules

    def disable_metric_groups(self, name, args):
        self.metric_groups = self.handler_modules[name]["metric_groups"]
        metrics_disable = list()
        if 'disable' not in self.metric_groups.keys():
            self.metric_groups.update({"disable": metrics_disable})

        for metric in args:
            metrics_disable.append(metric)
            if 'enable' in self.metric_groups.keys() and metric in self.metric_groups['enable']:
                self.metric_groups['enable'].remove(metric)

        self.metric_groups['disable'] = metrics_disable

        return self.handler_modules

    def remove_metric_groups(self, name, args):
        self.metric_groups = self.handler_modules[name]["metric_groups"]

        for metric in args:
            if 'enable' in self.metric_groups.keys() and metric in self.metric_groups['enable']:
                self.metric_groups['enable'].remove(metric)
            if 'disable' in self.metric_groups.keys() and metric in self.metric_groups['disable']:
                self.metric_groups['disable'].remove(metric)

        return self.handler_modules

    def remove_filters(self, name, *args):

        self.handler_modules[name]["filter"] = UtilsManager.remove_configs(self, self.handler_modules[name]["filter"],
                                                                           *args)

        return self.handler_modules

    def remove_configs(self, name, *args):

        self.handler_modules[name]["config"] = UtilsManager.remove_configs(self, self.handler_modules[name]["config"],
                                                                           *args)

        return self.handler_modules

    def remove_module(self, name):
        assert_that(name, is_in(list(self.handler_modules.keys())), "Invalid module")
        self.handler_modules.pop(name)
        return self.handler_modules

    def json(self):
        return json.dumps(self.handler_modules)


class Policy(HandlerModules, HandlerConfigs):
    def __init__(self, name, description, backend_type):

        self.policy = {"name": name,
                       "description": description,
                       "backend": backend_type,
                       "policy": {"handlers": {
                           "config": {},
                           "modules": {}
                       },
                           "input": {},
                           "config": {},
                           "kind": "collection"
                       }}
        self.config = self.policy['policy']['config']
        self.handler_configs = self.policy['policy']["handlers"]["config"]
        self.handler_modules = self.policy['policy']["handlers"]["modules"]

    def add_module_configs(self, name, **kwargs):
        self.handler_modules[name]['config'] = UtilsManager.add_configs(self, self.handler_modules[name]['config'],
                                                                        **kwargs)
        return self.policy

    def remove_module_configs(self, name, *args):
        self.handler_modules[name]['config'] = UtilsManager.remove_configs(self, self.handler_modules[name]['config'],
                                                                           *args)
        return self.policy

    def add_module_filters(self, name, **kwargs):
        self.handler_modules[name]['filter'] = UtilsManager.add_filters(self, self.handler_modules[name]['filter'],
                                                                        **kwargs)
        return self.policy

    def remove_module_filters(self, name, *args):
        self.handler_modules[name]['filter'] = UtilsManager.remove_filters(self, self.handler_modules[name]['filter'],
                                                                           *args)
        return self.policy

    def add_handler_configs(self, **kwargs):
        self.handler_configs = UtilsManager.add_configs(self, self.handler_configs, **kwargs)
        return self.policy

    def remove_handler_configs(self, *args):
        self.handler_configs = UtilsManager.remove_configs(self, self.handler_configs, *args)
        return self.policy

    def add_input_configs(self, **kwargs):
        assert_that('input_type', is_in(list(self.policy['policy']['input'].keys())),
                    "It is not possible to enter settings without defining the input. Use `add_input` first.")
        if 'tap' not in self.policy['policy']['input'].keys() and 'tap_selector' not in self.policy['policy'][
            'input'].keys():
            raise ValueError("It is not possible to enter settings without defining the input. Use `add_input` first")
        if 'config' not in self.policy['policy']['input'].keys():
            self.policy['policy']['input'].update({'config': {}})
        self.policy['policy']['input']['config'] = UtilsManager.add_configs(self,
                                                                            self.policy['policy']['input']['config'],
                                                                            **kwargs)
        return self.policy

    def remove_input_configs(self, *args):
        self.policy['policy']['input']['config'] = UtilsManager.remove_configs(self, self.policy['input']['config'],
                                                                               *args)
        return self.policy

    def add_input_filters(self, **kwargs):
        assert_that('input_type', is_in(list(self.policy['policy']['input'].keys())),
                    "It is not possible to enter settings without defining the input. Use `add_input` first.")
        if 'tap' not in self.policy['policy']['input'].keys() and 'tap_selector' not in self.policy['policy'][
            'input'].keys():
            raise ValueError("It is not possible to enter settings without defining the input. Use `add_input` first")
        if 'filter' not in self.policy['policy']['input'].keys():
            self.policy['policy']['input'].update({'filter': {}})
        self.policy['policy']['input']['filter'] = UtilsManager.add_filters(self,
                                                                            self.policy['policy']['input']['filter'],
                                                                            **kwargs)
        return self.policy

    def remove_input_filters(self, name, *args):
        self.policy['policy']['input']['filter'] = UtilsManager.remove_filters(self,
                                                                               self.policy['policy']['input']['filter'],
                                                                               *args)
        return self.policy

    def add_configs(self, **kwargs):
        self.config = UtilsManager.add_configs(self, self.config, **kwargs)
        return self.policy

    def remove_configs(self, *args):
        self.config = UtilsManager.remove_configs(self, self.config, *args)
        return self.policy

    def add_filters(self, **kwargs):
        raise ValueError(f"Policy objects do not have filters. Try `add_module_filters` or `add_input_filters` instead")

    def remove_filters(self, **kwargs):
        raise ValueError(
            f"Policy objects do not have filters. Try `remove_module_filters` or `remove_input_filters` instead")

    def __add_input_tap(self, input_type, name):
        assert_that('tap_selector', not_(is_in(list(self.policy['policy']['input'].keys()))),
                    "tap_selector is already defined. Use `remove_input` first.")
        if 'tap' not in self.policy['input'].keys():
            self.policy['policy']['input'].update({'tap': {}})
        if 'input_type' not in self.policy['policy']['input'].keys():
            self.policy['policy']['input'].update({'input_type': {}})
        self.policy['policy']['input']['tap'] = name
        self.policy['policy']['input']['input_type'] = input_type
        return self.policy

    def __add_input_tap_selector(self, input_type, **kwargs):
        assert_that('tap', not_(is_in(list(self.policy['policy']['input'].keys()))),
                    "tap is already defined. Use `remove_input` first.")
        assert_that('input_match', is_in(list(kwargs.keys())),
                    "`input_match` is a required parameter if selector is `tap_selector`")
        assert_that('tags', is_in(list(kwargs.keys())),
                    "`tags` is a required parameter if selector is `tap_selector`")
        assert_that(kwargs['input_match'], any_of(equal_to('any'), equal_to('all')), "Invalid input_match")
        input_match = kwargs['input_match']
        kwargs.pop('input_match')
        if 'tap_selector' not in self.policy['policy']['input'].keys():
            self.policy['policy']['input'].update({'tap_selector': {}})
        if 'input_type' not in self.policy['policy']['input'].keys():
            self.policy['policy']['input'].update({'input_type': {}})
            all_selectors = list()
        elif input_match in self.policy['policy']['input']['tap_selector'].keys():
            all_selectors = self.policy['policy']['input']['tap_selector'][input_match]
        else:
            all_selectors = list()

        for selector_key in kwargs['tags']:
            all_selectors.append({selector_key: kwargs['tags'][selector_key]})

        self.policy['policy']['input']['tap_selector'] = {input_match: all_selectors}
        self.policy['policy']['input']['input_type'] = input_type

    def add_input(self, input_type, selector, **kwargs):
        assert_that(selector, any_of('tap', 'tap_selector'), "Invalid input selector")

        if selector == 'tap':
            assert_that('name', is_in(list(kwargs.keys())),
                        "If `selector=tap`, you need to specify tap name. name=`the_name_you_want`.")
            self.__add_input_tap(input_type, kwargs['name'])

        else:
            assert_that('input_match', is_in(list(kwargs.keys())),
                        "If `selector=tap`, you need to specify input_match. input_match=`any` or input_match=`all`.")
            self.__add_input_tap_selector(input_type, **kwargs)

    def remove_input(self):
        self.policy['policy']['input'] = dict()

    def json(self):
        return json.dumps(self.policy)
