from behave import when, then
from test_config import TestConfig
from ui_utils import *
from hamcrest import *
from utils import random_string, create_tags_set
from page_objects import *
from control_plane_sink import sink_name_prefix
from selenium.common.exceptions import TimeoutException, StaleElementReferenceException

configs = TestConfig.configs()
orb_url = configs.get('orb_url')
sink_description_prefix = "sink_description_"
sink_remote_url = "www.remoteurl.com"
username = configs.get('email')
password = configs.get('password')
  
    
@when('a sink is created through the UI with {orb_tags} orb tag')
def create_sink(context, orb_tags):
    context.orb_tags = create_tags_set(orb_tags)
    context.initial_counter_datatable = check_total_counter(context.driver)
    WebDriverWait(context.driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, SinkPage.new_sink_button())),
        message="Unable to click on new sink"
                " button").click()
    WebDriverWait(context.driver, 5).until(EC.url_to_be(f"{orb_url}/pages/sinks/add"), 
                                           message="Orb add Sink Management page not available")
    context.name_label = sink_name_prefix + random_string(5)
    input_text_by_xpath(SinkPage.name_label(), context.name_label, context.driver)
    context.sink_description = sink_description_prefix + random_string(5)
    input_text_by_xpath(SinkPage.sink_description(), context.sink_description, context.driver)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()
    context.remote_url = sink_remote_url
    input_text_by_xpath(SinkPage.remote_url(), context.remote_url, context.driver)
    context.sink_username = username 
    input_text_by_xpath(SinkPage.sink_username(), context.sink_username, context.driver)
    context.sink_password = password
    input_text_by_xpath(SinkPage.sink_password(), context.sink_password, context.driver)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()    
    for tag_key, tag_value in context.orb_tags.items():
        input_text_by_xpath(SinkPage.sink_tag_key(), tag_key, context.driver)
        input_text_by_xpath(SinkPage.sink_tag_value(), tag_value, context.driver)
        WebDriverWait(context.driver, 3).until(
            EC.element_to_be_clickable((By.XPATH, SinkPage.sink_add_tag_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.save_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Sink successfully created'))
    context.initial_counter = check_total_counter(context.driver)
    
    
@then("the new sink {condition} shown on the datatable")
def check_presence_of_group_on_orb_ui(context, condition):
    sink = find_element_on_datatable_by_condition(context.driver, DataTable.sink_name_on_datatable(context.name_label),
                                                   LeftMenu.sinks(), condition)
    if condition == "is":
        assert_that(sink, is_not(None), "Unable to sink name on the datatable")
    else:
        assert_that(sink, is_(None), f"{context.name_label} found on group datable")    
        
@then("total number was increased in one unit")
def check_total_counter_final(context):
    final_counter_datatable = check_total_counter(context.driver)
    assert_that(final_counter_datatable, equal_to(context.initial_counter_datatable + 1),
                'The counter was increase with successfully')


def check_total_counter(driver):
    threading.Event().wait(3)
    WebDriverWait(driver, 3).until(
        EC.presence_of_element_located((By.XPATH, DataTable.page_count())))
    return int(driver.find_element(By.XPATH, DataTable.page_count()).text.split()[0])

@when("delete a sink using filter by name with {orb_tags} orb tag")
def delete_a_sink_item(context, orb_tags):
    create_sink(context, orb_tags)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.filter_by()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.filter_by())))
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.option_list())))
    select_list = WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.all_filter_options())))
    select_list[0].click()
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_element_located((By.XPATH, DataTable.filter_by_name_field())))
    input_text_by_xpath(DataTable.filter_by_name_field(), context.name_label, context.driver)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.plus_button()))).click()
    WebDriverWait(context.driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.trash_icon()))).click()
    input_text_by_xpath(SinkPage.delete_sink_confirmation_field(), context.name_label, context.driver)
    WebDriverWait(context.driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, SinkPage.delete_sink_confirmation_title()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, SinkPage.delete_sink_confirmation_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Sink successfully deleted'))
    WebDriverWait(context.driver, 3).until(EC.element_to_be_clickable((By.XPATH, UtilButton.clear_all_filters())),
                                           "Unable to clear all filters").click()



@when("delete a sink using filter by {filter_option} with {orb_tags} orb tag")
def delete_a_sink_item(context, filter_option, orb_tags):
    create_sink(context, orb_tags)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.filter_by()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.filter_by())))
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.option_list())))
    select_list = WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.all_filter_options())))
    select_list[0].click()
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_element_located((By.XPATH, DataTable.filter_by_name_field())))
    input_text_by_xpath(DataTable.filter_by_name_field(), context.name_label, context.driver)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.plus_button()))).click()
    WebDriverWait(context.driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.trash_icon()))).click()
    input_text_by_xpath(SinkPage.delete_sink_confirmation_field(), context.name_label, context.driver)
    WebDriverWait(context.driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, SinkPage.delete_sink_confirmation_title()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, SinkPage.delete_sink_confirmation_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Sink successfully deleted'))
    WebDriverWait(context.driver, 3).until(EC.element_to_be_clickable((By.XPATH, UtilButton.clear_all_filters())),
                                           "Unable to clear all filters").click()
    

@when("update a sink using filter by name with {orb_tags} orb tag")  
def update_a_sink_item(context, orb_tags):
    create_sink(context, orb_tags)
    context.initial_counter_datatable = check_total_counter(context.driver)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.filter_by()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.filter_by())))
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.option_list())))
    select_list = WebDriverWait(context.driver, 3).until(
        EC.presence_of_all_elements_located((By.XPATH, DataTable.all_filter_options())))
    select_list[0].click()
    WebDriverWait(context.driver, 3).until(
        EC.presence_of_element_located((By.XPATH, DataTable.filter_by_name_field())))
    input_text_by_xpath(DataTable.filter_by_name_field(), context.name_label, context.driver)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.plus_button()))).click()
    WebDriverWait(context.driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, DataTable.edit_icon()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, SinkPage.sink_description()))).click()
    WebDriverWait(context.driver, 5).until(
        EC.element_to_be_clickable((By.XPATH, SinkPage.name_label()))).clear()
    context.name_label = sink_name_prefix + "upd" + random_string(5)
    input_text_by_xpath(SinkPage.name_label(), context.name_label, context.driver)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()
    context.sink_password = password
    input_text_by_xpath(SinkPage.sink_password(), context.sink_password, context.driver)
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click() 
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.next_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.element_to_be_clickable((By.XPATH, UtilButton.save_button()))).click()
    WebDriverWait(context.driver, 3).until(
        EC.text_to_be_present_in_element((By.CSS_SELECTOR, "span.title"), 'Sink successfully updated'))
    WebDriverWait(context.driver, 3).until(EC.element_to_be_clickable((By.XPATH, UtilButton.clear_all_filters())),
                                           "Unable to clear all filters").click()
    context.initial_counter = check_total_counter(context.driver)
    

@then("total number was the same")
def check_total_counter_final(context):
    final_counter_datatable = check_total_counter(context.driver)
    assert_that(final_counter_datatable, equal_to(context.initial_counter_datatable),
                'The counter was not the same') 
