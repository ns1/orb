import re
from configs import TestConfig
from local_agent import get_orb_agent_logs
from users import get_auth_token
from utils import (random_string, filter_list_by_parameter_start_with, generate_random_string_with_predefined_prefix, \
                   create_tags_set, check_logs_contain_message_and_name, threading_wait_until, validate_json,
                   return_api_post_response, return_api_get_response, return_api_delete_response,
                   return_api_put_response)
from behave import given, then, step
from hamcrest import *
import requests
from random import sample
import json
import random

configs = TestConfig.configs()
agent_group_name_prefix = 'test_group_name_'
agent_group_description = "This is an agent group"
orb_url = configs.get('orb_url')
verify_ssl_bool = eval(configs.get('verify_ssl').title())


@step("{amount_of_agent_groups} Agent Group(s) is created with {amount_of_tags} tags contained in the agent")
def create_agent_group_matching_agent(context, amount_of_agent_groups, amount_of_tags, **kwargs):
    if amount_of_tags.isdigit() is False:
        assert_that(amount_of_tags, equal_to("all"), 'Unexpected value for amount of tags')

    if "group_description" in kwargs.keys():
        group_description = kwargs["group_description"]
    else:
        group_description = agent_group_description

    tags_in_agent = context.agent["orb_tags"]
    if context.agent["agent_tags"] is not None:
        tags_in_agent.update(context.agent["agent_tags"])
    tags_keys = tags_in_agent.keys()

    if amount_of_tags.isdigit() is True:
        amount_of_tags = int(amount_of_tags)
    else:
        amount_of_tags = len(tags_keys)
    assert_that(tags_keys, has_length(greater_than_or_equal_to(amount_of_tags)), "Amount of tags greater than tags"
                                                                                 "contained in agent")
    tags_to_group = {key: tags_in_agent[key] for key in sample(tags_keys, amount_of_tags)}
    assert_that(len(tags_to_group), greater_than(0), f"Unable to create group without tags. Tags:{tags_to_group}. "
                                                     f"Agent:{context.agent}")
    for group in range(int(amount_of_agent_groups)):
        agent_group_name = agent_group_name_prefix + random_string()
        agent_group_data = generate_group_with_valid_json(context.token, agent_group_name, group_description,
                                                          tags_to_group, context.agent_groups)


@step("{amount_of_agent_groups} Agent Group(s) is created with {orb_tags} orb tag(s) (lower case)")
# this step is temporary because of issue https://github.com/orb-community/orb/issues/1053
def create_group_with_tags_lower_case(context, amount_of_agent_groups, orb_tags):
    create_new_agent_group(context, amount_of_agent_groups, orb_tags, tags_lower_case=True)


@step("{amount_of_agent_groups} Agent Group(s) is created with {orb_tags} orb tag(s)")
def create_new_agent_group(context, amount_of_agent_groups, orb_tags, **kwargs):
    if "group_description" in kwargs.keys():
        group_description = kwargs["group_description"]
    else:
        group_description = agent_group_description
    if "tags_lower_case" in kwargs.keys() and kwargs["tags_lower_case"] is True:
        orb_tags = create_tags_set(orb_tags)
        context.orb_tags = {k.lower(): v.lower() for k, v in orb_tags.items()}
    else:
        context.orb_tags = create_tags_set(orb_tags)

    for group in range(int(amount_of_agent_groups)):
        agent_group_name = generate_random_string_with_predefined_prefix(agent_group_name_prefix)
        if len(context.orb_tags) == 0:
            context.agent_group_data = create_agent_group(context.token, agent_group_name, group_description,
                                                          context.orb_tags, 400)
        else:
            agent_group_data = generate_group_with_valid_json(context.token, agent_group_name, group_description,
                                                              context.orb_tags, context.agent_groups)
            group_id = agent_group_data['id']
            context.agent_groups[group_id] = agent_group_name


@step("{amount_of_agent_groups} Agent Group(s) is created with {orb_tags} orb tag(s) and {description} description")
def create_new_agent_group_with_defined_description(context, amount_of_agent_groups, orb_tags, description):
    for group in range(int(amount_of_agent_groups)):
        if description == "without":
            create_new_agent_group(context, amount_of_agent_groups, orb_tags, group_description=None)
        else:
            description = description.replace('"', '')
            description = description.replace(' as', '')
            create_new_agent_group(context, amount_of_agent_groups, orb_tags, group_description=description)


@step("{amount_of_agent_groups} Agent Group(s) is created with same tag as the agent and {description} description")
def create_agent_group_with_defined_description_and_matching_agent(context, amount_of_agent_groups, description):
    if description == "without":
        create_agent_group_matching_agent(context, amount_of_agent_groups, "all", group_description=None)
    else:
        if description == "with":
            description = f"test agent group description {random_string()}"
        else:
            description = description.replace('"', '')
            description = description.replace(' as', '')
        create_agent_group_matching_agent(context, amount_of_agent_groups, "all", group_description=description)


@step("{group_order} agent group {edited_parameter} must be empty")
def check_if_value_is_empty_after_editing(context, group_order, edited_parameter):
    agent_group_id = get_group_by_order(group_order, list(context.agent_groups.keys()))
    group_after_editing = get_agent_group(context.token, agent_group_id)
    if edited_parameter in group_after_editing.keys():
        assert_that(group_after_editing[edited_parameter], equal_to(any_of(None, "", {})),
                    f"Agent group {edited_parameter} must be empty, but is not. "
                    f"Group after editing: {group_after_editing}")
    else:
        assert_that(edited_parameter, is_in(context.group_before_editing.keys()), f"{edited_parameter} "
                                                                                  f" was not already present in group "
                                                                                  f"before editing. Before editing: "
                                                                                  f"{context.group_before_editing}\n. "
                                                                                  f"After Editing: "
                                                                                  f"{group_after_editing}\n")
        assert_that(edited_parameter, not_(is_in(group_after_editing.keys())), f"{edited_parameter} must be empty, but "
                                                                               f"is present in group with non empty "
                                                                               f"values after editing. Before "
                                                                               f"editing: "
                                                                               f"{context.group_before_editing}\n. "
                                                                               f"After Editing: "
                                                                               f"{group_after_editing}\n")


@step("{group_order} agent group {edited_parameter} must remain the same")
def check_if_value_is_the_same_as_before(context, group_order, edited_parameter):
    agent_group_id = get_group_by_order(group_order, list(context.agent_groups.keys()))
    group_after_editing = get_agent_group(context.token, agent_group_id)
    if edited_parameter in context.group_before_editing.keys() and edited_parameter in group_after_editing.keys():
        assert_that(context.group_before_editing[edited_parameter], equal_to(group_after_editing[edited_parameter]),
                    f"Agent group {edited_parameter} has different value than before editing. Group before editing: "
                    f"{context.group_before_editing}. Group after editing: {group_after_editing}")
    else:
        assert_that(edited_parameter, is_in(context.group_before_editing.keys()), f"{edited_parameter} "
                                                                                  f"not present in group before "
                                                                                  f"editing. Before editing: "
                                                                                  f"{context.group_before_editing}\n. "
                                                                                  f"After Editing: "
                                                                                  f"{group_after_editing}\n")
        assert_that(edited_parameter, is_in(group_after_editing.keys()), f"{edited_parameter} not present in group "
                                                                         f"after editing. Before editing: "
                                                                         f"{context.group_before_editing}\n. "
                                                                         f"After Editing: "
                                                                         f"{group_after_editing}\n")


@step("the {edited_parameters} of {group_order} Agent Group is edited using: {parameters_values}")
def edit_multiple_groups_parameters(context, edited_parameters, group_order, parameters_values):
    edited_parameters = edited_parameters.split(", ")
    assert_that(group_order, any_of(equal_to("first"), equal_to("second"), equal_to("last")),
                "Unexpected value for group.")
    order_convert = {"first": 0, "last": -1, "second": 1}
    agent_groups_id = list(context.agent_groups.keys())[order_convert[group_order]]
    for param in edited_parameters:
        assert_that(param, any_of(equal_to('name'), equal_to('description'), equal_to('tags')),
                    'Unexpected parameter to edit')
    parameters_values = parameters_values.split("/ ")

    group_editing = get_agent_group(context.token, agent_groups_id)
    group_data = {"name": group_editing["name"], "tags": group_editing["tags"]}
    if "description" in group_editing.keys():
        group_data["description"] = group_editing["description"]
    else:
        group_data["description"] = None

    editing_param_dict = dict()
    for param in parameters_values:
        param_split = param.split("=")
        if param_split[1].lower() == "empty":
            param_split[1] = {}
        elif param_split[1].lower() == "omitted":
            param_split[1] = None
        editing_param_dict[param_split[0]] = param_split[1]

    assert_that(set(editing_param_dict.keys()), equal_to(set(edited_parameters)),
                "All parameter must have referenced value")

    if "tags" in editing_param_dict.keys() and editing_param_dict["tags"] is not None and editing_param_dict["tags"] \
            != {}:
        if re.match(r"matching (\d+|all|the) agent*", editing_param_dict["tags"]):
            # todo improve logic for multiple agents
            editing_param_dict["tags"] = context.agent["orb_tags"]
            if context.agent["agent_tags"] is not None:
                editing_param_dict["tags"].update(context.agent["agent_tags"])
        else:
            editing_param_dict["tags"] = create_tags_set(editing_param_dict["tags"])
    expected_status_code = 200
    if "name" in editing_param_dict.keys() and editing_param_dict["name"] is not None and editing_param_dict["name"] \
            != {}:
        if editing_param_dict['name'] == "conflict":
            agent_group_name = list(context.agent_groups.values())[-1]
            editing_param_dict["name"] = agent_group_name
            expected_status_code = 409
        else:
            editing_param_dict["name"] = f"{agent_group_name_prefix}{editing_param_dict['name']}_{random_string(5)}"

    for parameter, value in editing_param_dict.items():
        group_data[parameter] = value

    context.group_before_editing = get_agent_group(context.token, agent_groups_id)
    context.editing_response, real_status_code = edit_agent_group(context.token, agent_groups_id, group_data["name"],
                                                                  group_data["description"], group_data["tags"],
                                                                  expected_status_code=expected_status_code)
    if real_status_code >= 300 or real_status_code < 200:
        context.error_message = context.editing_response


@then("agent group editing must fail")
def fail_group_editing(context):
    assert_that(list(context.editing_response.keys())[0], equal_to("error"), f"Agent group editing process was supposed"
                                                                             f" to fail, but didn't.")


@step("Agent Group creation response must be an error with message '{message}'")
def error_response_message(context, message):
    response = list(context.agent_group_data.items())[0]
    response_key, response_value = response[0], response[1]
    assert_that(response_key, equal_to('error'),
                'Response of invalid agent group creation must be an error')
    assert_that(response_value, equal_to(message), "Unexpected message for error")


@step("{amount_agent_matching} agent must be matching on response field matching_agents of the {group_order} group"
      " created")
def matching_agent(context, amount_agent_matching, group_order):
    assert_that(group_order, any_of(equal_to("first"), equal_to("second"), equal_to("last")),
                "Unexpected value for group.")
    order_convert = {"first": 0, "last": -1, "second": 1}
    agent_groups_id = list(context.agent_groups.keys())[order_convert[group_order]]
    agent_group_data = get_agent_group(context.token, agent_groups_id)
    matching_total_agents = agent_group_data['matching_agents']['total']
    assert_that(matching_total_agents, equal_to(int(amount_agent_matching)))


@step("{amount_of_groups_to_remove} group(s) to which the agent is linked is removed")
def remove_group(context, amount_of_groups_to_remove):
    container_logs = get_orb_agent_logs(context.container_id)
    amount_of_groups_to_remove = int(amount_of_groups_to_remove)
    assert_that(len(list(context.agent['last_hb_data']['group_state'].keys())),
                greater_than_or_equal_to(amount_of_groups_to_remove),
                f"The number of groups to be removed cannot be greater than the number to which the agent is subscribed"
                f"\nAgent: {context.agent}."
                f"\nAgent logs:{container_logs}")
    group_linked_to_remove_id = random.sample(list(context.agent['last_hb_data']['group_state'].keys()),
                                              amount_of_groups_to_remove)
    for group in group_linked_to_remove_id:
        delete_agent_group(context.token, group)
        context.agent_groups.pop(group)
    context.removed_groups_ids = group_linked_to_remove_id


@then('cleanup agent group')
def clean_agent_groups(context):
    """
    Remove all agent groups starting with 'agent_group_name_prefix' from the orb

    :param context: Behave object that contains contextual information during the running of tests.
    """
    token = context.token
    agent_groups_list = list_agent_groups(token)
    agent_groups_filtered_list = filter_list_by_parameter_start_with(agent_groups_list, 'name', agent_group_name_prefix)
    delete_agent_groups(token, agent_groups_filtered_list)


@step("referred agent is subscribed to {amount_of_groups} {group}")
def subscribe_agent_to_a_group(context, amount_of_groups, group):
    assert_that(group, any_of(equal_to("group"), equal_to("groups")), "Unexpected word on step description")
    agent = context.agent
    agent_tags = agent['orb_tags']
    if agent["agent_tags"] is not None:
        agent_tags.update(agent["agent_tags"])
    for group in range(int(amount_of_groups)):
        agent_group_name = generate_random_string_with_predefined_prefix(agent_group_name_prefix)
        agent_group_data = generate_group_with_valid_json(context.token, agent_group_name,
                                                          agent_group_description, agent_tags,
                                                          context.agent_groups)
        assert_that(agent_group_data['matching_agents']['online'], equal_to(1),
                    f"No agent matching this group.\n\n {agent_group_data}. \n"
                    f"Agent: {agent}")


@step('the container logs contain the message "{text_to_match}" referred to each matching group within'
      '{time_to_wait} seconds')
def check_logs_for_group(context, text_to_match, time_to_wait):
    groups_matching, context.groups_matching_id = return_matching_groups(context.token, context.agent_groups,
                                                                         context.agent)
    text_found, groups_to_which_subscribed = check_subscription(groups_matching, text_to_match, context.container_id,
                                                                timeout=time_to_wait)
    container_logs = get_orb_agent_logs(context.container_id)
    assert_that(text_found, is_(True), f"Message {text_to_match} was not found in the agent logs for group(s)"
                                       f"{set(groups_matching).difference(groups_to_which_subscribed)}!.\n\n"
                                       f"Logs = {container_logs}. \n\n"
                                       f"Agent: {json.dumps(context.agent, indent=4)} \n\n")


@step("a new group is requested to be created with the same name as an existent one")
def create_group_with_name_conflict(context):
    tags = create_tags_set('1')
    name = list(context.agent_groups.values())[0]
    context.error_message = create_agent_group(context.token, name, 'conflict_group', tags, 409)


def create_agent_group(token, name, description, tags, expected_status_code=201):
    """
    Creates an agent group in Orb control plane

    :param (str) token: used for API authentication
    :param (str) name: of the agent group to be created
    :param (str) description: description of group
    :param (dict) tags: dict with all pairs key:value that will be used as tags
    :returns: (dict) a dictionary containing the created agent group data
    :param (int) expected_status_code: expected request's status code. Default:201 (happy path).
    """

    json_request = {"name": name, "description": description, "tags": tags}
    if expected_status_code == 201:
        assert_that(len(tags), greater_than(0), f"Tags is required to created a group. Json used: {json_request}")
    status_code, response = return_api_post_response(f"{orb_url}/api/v1/agent_groups",
                                                     request_body=json_request, token=token, verify=verify_ssl_bool)

    assert_that(status_code, equal_to(expected_status_code),
                f"Request to create agent group failed with status= {str(status_code)}. Response="
                f" {str(response)}. Json used: {json_request}")

    return response


def get_agent_group(token, agent_group_id):
    """
    Gets an agent group from Orb control plane

    :param (str) token: used for API authentication
    :param (str) agent_group_id: that identifies the agent group to be fetched
    :returns: (dict) the fetched agent group
    """
    status_code, response = return_api_get_response(f"{orb_url}/api/v1/agent_groups/{agent_group_id}",
                                                    token=token,
                                                    verify=verify_ssl_bool)
    assert_that(status_code, equal_to(200),
                'Request to get agent group id=' + agent_group_id + ' failed with status=' + str(
                    status_code))

    return response


def list_agent_groups(token, limit=100, offset=0):
    """
    Lists all agent_groups from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :param (int) limit: Size of the subset to retrieve. (max 100). Default = 100
    :param (int) offset: Number of items to skip during retrieval. Default = 0.
    :returns: (list) a list of agent_groups
    """
    all_agent_groups, total, offset = list_up_to_limit_agent_groups(token, limit, offset)

    new_offset = limit + offset

    while new_offset < total:
        agent_groups_from_offset, total, offset = list_up_to_limit_agent_groups(token, limit, new_offset)
        all_agent_groups = all_agent_groups + agent_groups_from_offset
        new_offset = limit + offset

    return all_agent_groups


def list_up_to_limit_agent_groups(token, limit=100, offset=0):
    """
    Lists up to 100 agent groups from Orb control plane that belong to this user

    :param (str) token: used for API authentication
    :param (int) limit: Size of the subset to retrieve (max 100). Default = 100
    :param (int) offset: Number of items to skip during retrieval. Default = 0.
    :returns: (list) a list of agent groups, (int) total groups on orb, (int) offset
    """
    status_code, response = return_api_get_response(f"{orb_url}/api/v1/agent_groups",
                                                    token=token, verify=verify_ssl_bool,
                                                    params={"limit": limit, "offset": offset})

    assert_that(status_code, equal_to(200),
                'Request to list agent groups failed with status=' + str(status_code))

    assert_that(response, has_key('agentGroups'), f"Response does not contain agentGroups. Response: {str(response)}")
    assert_that(response, has_key('total'), f"Response does not contain total. Response: {str(response)}")
    assert_that(response, has_key('offset'), f"Response does not contain offset. Response: {str(response)}")

    return response['agentGroups'], response['total'], response['offset']


def delete_agent_groups(token, list_of_agent_groups):
    """
    Deletes from Orb control plane the agent groups specified on the given list

    :param (str) token: used for API authentication
    :param (list) list_of_agent_groups: that will be deleted
    """

    for agent_Groups in list_of_agent_groups:
        delete_agent_group(token, agent_Groups['id'])


def delete_agent_group(token, agent_group_id):
    """
    Deletes an agent group from Orb control plane

    :param (str) token: used for API authentication
    :param (str) agent_group_id: that identifies the agent group to be deleted
    """
    status_code, response = return_api_delete_response(f"{orb_url}/api/v1/agent_groups/{agent_group_id}",
                                                       token=token, verify=verify_ssl_bool)
    assert_that(status_code, equal_to(204), 'Request to delete agent group id='
                + agent_group_id + ' failed with status=' + str(status_code))


@threading_wait_until
def check_subscription(agent_groups_names, expected_message, container_id, event=None):
    """

    :param (list) agent_groups_names: groups to which the agent must be subscribed
    :param (str) expected_message: message that we expect to find in the logs
    :param (str) container_id: agent container id
    :param (obj) event: threading.event
    :return: (bool) True if agent is subscribed to all matching groups, (list) names of the groups to which agent is subscribed
    """
    groups_to_which_subscribed = set()
    for name in agent_groups_names:
        logs = get_orb_agent_logs(container_id)
        text_found, log_line = check_logs_contain_message_and_name(logs, expected_message, name, "group_name")
        if text_found is True:
            groups_to_which_subscribed.add(log_line["group_name"])
            if set(groups_to_which_subscribed) == set(agent_groups_names):
                event.set()
                return event.is_set(), groups_to_which_subscribed

    return event.is_set(), groups_to_which_subscribed


def edit_agent_group(token, agent_group_id, name, description, tags, expected_status_code=200):
    """

    :param (str) token: used for API authentication
    :param (str) agent_group_id: that identifies the agent group to be edited
    :param (str) name: agent group's name
    :param (str) description: agent group's description
    :param (str) tags: orb tags that will be used to connect agents to groups
    :param (int) expected_status_code: expected request's status code. Default:200.
    :returns: (dict) the edited agent group
    """

    if description == {}:
        description = ""

    json_request = {"name": name, "description": description, "tags": tags,
                    "validate_only": False}
    json_request = {parameter: value for parameter, value in json_request.items() if value is not None}

    status_code, response = return_api_put_response(f"{orb_url}/api/v1/agent_groups/{agent_group_id}",
                                                    request_body=json_request, token=token, verify=verify_ssl_bool)
    if tags == {} or name == {}:
        expected_status_code = 400

    assert_that(status_code, equal_to(expected_status_code),
                'Request to edit agent group failed with status=' + "status code =" +
                str(status_code) + "response =" + str(response) +
                " json used: " + str(json_request))

    return response, status_code


def return_matching_groups(token, existing_agent_groups, agent_json):
    """

    :param (str) token: used for API authentication
    :param (dict) existing_agent_groups: dictionary with the existing groups, the id of the groups being the key and the name the values
    :param (dict) agent_json: dictionary containing all the information of the agent to which the groups must be matching

    :return (list): groups_matching, groups_matching_id
    """
    groups_matching = list()
    groups_matching_id = list()
    for group in existing_agent_groups.keys():
        group_data = get_agent_group(token, group)
        group_tags = dict(group_data["tags"])
        agent_tags = agent_json["orb_tags"]
        agent_tags.update(agent_json['agent_tags'])
        if all(item in agent_tags.items() for item in group_tags.items()) is True:
            groups_matching.append(existing_agent_groups[group])
            groups_matching_id.append(group)
    return groups_matching, groups_matching_id


def tags_to_match_k_groups(token, k, all_existing_groups):
    """

    :param (str) token: used for API authentication
    :param (str) k: amount of groups that the agent must match
    :param (dict) all_existing_groups: full data of all existing groups
    :return: (dict) with all tags that must be on the agent
    """
    if k.isdigit() is False:
        assert_that(k, any_of(equal_to("all"), equal_to("last"), equal_to("first")),
                    "Unexpected amount of groups to match.")
    if k == "last":
        id_of_groups_to_match = [list(all_existing_groups.keys())[-1]]
    elif k == "first":
        id_of_groups_to_match = [list(all_existing_groups.keys())[0]]
    else:
        if k == "all":
            k = len(list(all_existing_groups.keys()))
        id_of_groups_to_match = sample(list(all_existing_groups.keys()), int(k))
    all_used_tags = dict()
    for agent_group_id in id_of_groups_to_match:
        group_data = get_agent_group(token, agent_group_id)
        all_used_tags.update(group_data["tags"])
    return all_used_tags


def generate_group_with_valid_json(token, agent_group_name, group_description, tags_to_group, agent_groups):
    """
    Create a group and validate the json schema

    :param (str) token: used for API authentication
    :param (str) agent_group_name: of the agent group to be created
    :param (str) group_description: description of group
    :param (dict) tags_to_group: dict with all pairs key:value that will be used as tags
    :returns: (dict) a dictionary containing the created agent group data

    :return: agent group data
    """
    agent_group_data = create_agent_group(token, agent_group_name, group_description,
                                          tags_to_group)
    group_id = agent_group_data['id']
    agent_groups[group_id] = agent_group_name

    local_orb_path = configs.get("local_orb_path")
    agent_group_schema_path = local_orb_path + "/python-test/features/steps/schemas/groups_schema.json"
    is_schema_valid = validate_json(agent_group_data, agent_group_schema_path)
    assert_that(is_schema_valid, equal_to(True), f"Invalid group json. \n Group = {agent_group_data}")
    return agent_group_data


def get_group_by_order(group_order, agent_groups_ids):
    """

    :param group_order: first, second or last group created
    :param agent_groups_ids: list of all agent groups ids created on test process
    :return: id of group
    """
    assert_that(group_order, any_of(equal_to("first"), equal_to("second"), equal_to("last")),
                "Unexpected value for group.")
    order_convert = {"first": 0, "last": -1, "second": 1}
    agent_group_id = agent_groups_ids[order_convert[group_order]]
    return agent_group_id
