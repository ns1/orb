/**
 * Agent Policy Data Model
 *
 * [Agent Policy Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Policies-and-Datasets}
 */

import { OrbEntity } from 'app/common/interfaces/orb/orb.entity.interface';

/**
 * @interface AgentPolicy
 */
export interface AgentPolicy extends OrbEntity {
  /**
   * Description {string}
   */
  description?: string;

  /**
   * Agent Backend this policy is for: {string}
   * Cannot change once created (read only)
   */
  backend?: string;

  /**
   * Version
   * monotonically increases on each update
   * starts at 0
   */
  version?: number;

  /**
   * A timestamp of creation {string}
   */
  ts_created?: string;

  /**
   * Tags {{[propName: string]: string}}
   * User defined tags
   */
  tags?: any;

  /**
   * Policy Metadata {{[propName: string]: string}}
   */
  policy_metadata?: any;

}
