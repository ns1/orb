/**
 * Prometheus Sink Config Interface
 *
 * [Prometheus Sink]{@link https://github.com/ns1labs/orb/blob/develop/cmd/prom-sink/main.go}
 * [Sinks Architecture]{@link https://github.com/ns1labs/orb/wiki/Architecture:-Sinks}
 */
import { SinkConfig } from 'app/common/interfaces/orb/sink.config/sink.config.interface';

/**
 * @interface PrometheusConfig
 */
export interface PrometheusConfig extends SinkConfig<string> {

  /**
   *  Remote Host URL {string}
   */
  remote_host?: string;

  /**
   *  Username|Email(?) {string}
   */
  username?: string;

  /**
   *  Password {string}
   */
  password?: string;
}
