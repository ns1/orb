import { Observable } from 'rxjs';

export enum FilterTypes {
  Input, // string input
  AutoComplete,
  Select, // allows select one option
  MultiSelect, // allows select multi options
  Checkbox, // on|off option
  Number, // number input
}

export interface FilterOption {
  name: string;
  prop: string;
  filter: (item: any, prop: any, value: any, extra?: any) => any;
  type: FilterTypes;
  param?: any;
  options?: string[];
  exact?: any;
  autoSuggestion?: Observable<string[]>;
}

export function filterExact(item: string, value: string): boolean {
  return item === value;
}

export function filterSubstr(item: string, value: string) {
  return item.toLowerCase().includes(value.toLowerCase());
}

export function filterString(item: any, prop: any, value: any, exact?: any) {
  const itemValue = !!item[prop] && item[prop] || '';
  if (exact) {
    return filterExact(itemValue, value);
  } else {
    return filterSubstr(itemValue, value);
  }
}

export function filterNumber(item: any, prop: any, value: any) {
  const itemValue = !!item[prop] && item[prop] || 0;
  return itemValue.toString() === value;
}

export function filterTags(item: any, prop: any, value: any, exact?: any) {
  if (!item || typeof item[prop] !== 'object') {
    return false;
  }
  // tag values
  const values = Object.entries(item[prop]).map(
    (entry) => `${entry[0]}:${entry[1]}`,
  );
  const filterFn = exact ? filterExact : filterSubstr;
  const filterValue = value.replace(' ', '');
  return values.reduce((acc, val) => {
    acc = acc || filterFn(val, filterValue);
    return acc;
  }, false);
}

export function filterMultiSelect(item: any, prop: any, values: any, exact?: any) {
  return values.reduce((prev, cur) => {
    if (exact) {
      return item[prop] === cur || prev;
    } else {
      return item[prop].includes(cur) || prev;
    }
  }, false);
}
