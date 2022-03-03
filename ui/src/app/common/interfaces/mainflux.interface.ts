export interface UserGroup {
  id?: string;

  name?: string;

  description?: string;

  metadata?: Object;

  owner_id?: string;

  parent_id?: string;

  users?: User[];
}

export interface User {
  id?: string;

  email?: string;

  password?: string;

  picture?: string;

  metadata?: Object;
}

export interface Channel {
  id?: string;

  name?: string;

  metadata?: any;

  type?: string;
}

export interface Thing {
  id?: string;

  key?: string;

  name?: string;

  metadata?: any;

  type?: string;
}

export interface Attribute {
  name: string;

  channel: string;

  subtopic?: string;

  persist_state: boolean;
}

export interface Definition {
  id: string;

  created: Date;

  attributes: Attribute[];

  delta: number;
}

export interface Twin {
  name?: string;

  owner?: string;

  id?: string;

  revision?: number;

  created?: Date;

  updated?: Date;

  definitions?: Definition[];

  definition?: any; // for request
  metadata?: any;
}

export interface SenMLRec {
  bn: string;

  bt: number;

  bu: string;

  bver: number;

  n: string;

  t: number;

  u: string;

  v: number;

  vd: number;

  vb: number;

  vs: string;
}

export interface MainfluxMsg {
  name: string;

  time: number;

  unit: string;

  value: number;

  string_value: string;

  bool_value: boolean;

  data_value: string;

  sum: number;

  protocol: string;
}

export interface MsgFilters {
  offset?: number;

  limit?: number;

  publisher?: string;

  subtopic?: string;

  name?: string;

  v?: string;

  vs?: string;

  vd?: string;

  vb?: string;

  from?: number;

  to?: number;
}

export interface ReaderUrl {
  prefix?: string;

  suffix?: string;
}

export interface Message {
  value: number;

  time: number;
}

export interface Dataset {
  label?: string;

  messages?: Message[];
}

// Mainflux::ORB
export interface PageFilters {
  limit?: number;

  offset?: number;

  name?: string;

  order?: string;

  dir?: string;

  // mainflux and other components metadata
  metadata?: string;

  type?: string;

  // orb
  tags?: Object;
}

export interface DropdownFilterItem {
  id?: string;

  label?: string;

  selected?: boolean;

  prop?: string;

  filter?: (item: any, filter: any) => boolean;
}

export interface TableConfig {
  keys?: string[];

  colNames?: string[];
}

export interface TablePage {
  limit?: number;

  offset?: number;

  total?: number;

  rows?: Object[];
}
