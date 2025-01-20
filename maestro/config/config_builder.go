package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/orb-community/orb/pkg/errors"
	"gopkg.in/yaml.v2"
)

var k8sOtelCollector = `
{
  "kind": "List",
  "apiVersion": "v1",
  "metadata": {
  },
  "items": [
    {
      "kind": "ConfigMap",
      "apiVersion": "v1",
      "metadata": {
        "name": "otel-collector-config-SINK_ID",
        "creationTimestamp": null
      },
      "data": {
        "config.yaml": "SINK_CONFIG"
      }
    },
    {
      "kind": "Deployment",
      "apiVersion": "apps/v1",
      "metadata": {
        "name": "otel-SINK_ID",
        "creationTimestamp": null,
        "labels": {
          "app": "opentelemetry",
          "component": "otel-collector"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "app": "opentelemetry",
            "component": "otel-collector-SINK_ID"
          }
        },
        "template": {
          "metadata": {
            "annotations": {
              "prometheus.io/path": "/metrics",
              "prometheus.io/port": "8888",
              "prometheus.io/scrape": "true"
            },
            "creationTimestamp": null,
            "labels": {
              "app": "opentelemetry",
              "component": "otel-collector-SINK_ID"
            }
          },
          "spec": {
            "volumes": [
              {
                "name": "varlog",
                "hostPath": {
                  "path": "/var/log",
                  "type": ""
                }
              },
              {
                "name": "varlibdockercontainers",
                "hostPath": {
                  "path": "/var/lib/docker/containers",
                  "type": ""
                }
              },
              {
                "name": "data",
                "configMap": {
                  "name": "otel-collector-config-SINK_ID",
                  "defaultMode": 420
                }
              }
            ],
            "containers": [
              {
                "name": "otel-collector",
                "image": "orb-community/otelcol-orb-maestro@latest",
                "ports": [
                  {
                    "containerPort": 13133,
                    "protocol": "TCP"
                  },
                  {
                    "containerPort": 8888,
                    "protocol": "TCP"
                  }
                ],
                "volumeMounts": [
                  {
                    "name": "varlog",
                    "readOnly": true,
                    "mountPath": "/var/log"
                  },
                  {
                    "name": "varlibdockercontainers",
                    "readOnly": true,
                    "mountPath": "/var/lib/docker/containers"
                  },
                  {
                    "name": "data",
                    "readOnly": true,
                    "mountPath": "/etc/otelcol-contrib/config.yaml",
                    "subPath": "config.yaml"
                  }
                ],
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "imagePullPolicy": "IfNotPresent"
              }
            ],
            "restartPolicy": "Always",
            "terminationGracePeriodSeconds": 30,
            "dnsPolicy": "ClusterFirst",
            "securityContext": {
            },
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": "25%",
            "maxSurge": "25%"
          }
        },
        "revisionHistoryLimit": 10,
        "progressDeadlineSeconds": 600
      },
      "status": {
      }
    },
    {
      "kind": "Service",
      "apiVersion": "v1",
      "metadata": {
        "name": "otel-SINK_ID",
        "creationTimestamp": null,
        "labels": {
          "app": "opentelemetry",
          "component": "otel-collector-SINK_ID"
        }
      },
      "spec": {
        "ports": [
          {
            "name": "metrics",
            "protocol": "TCP",
            "port": 8888,
            "targetPort": 8888
          },
          {
            "name": "healthcheck",
            "protocol": "TCP",
            "port": 13133,
            "targetPort": 13133
          }
        ],
        "selector": {
          "component": "otel-collector-SINK_ID"
        },
        "type": "ClusterIP",
        "sessionAffinity": "None"
      },
      "status": {
        "loadBalancer": {
        }
      }
    }
  ]
}
`

var JsonService = `
{
  "kind": "Service",
  "apiVersion": "v1",
  "metadata": {
	"name": "otel-SINK_ID",
	"creationTimestamp": null,
	"labels": {
	  "app": "opentelemetry",
	  "component": "otel-collector-SINK_ID"
	}
  },
  "spec": {
	"ports": [
	  {
		"name": "metrics",
		"protocol": "TCP",
		"port": 8888,
		"targetPort": 8888
	  },
	  {
		"name": "healthcheck",
		"protocol": "TCP",
		"port": 13133,
		"targetPort": 13133
	  }
	],
	"selector": {
	  "component": "otel-collector-SINK_ID"
	},
	"type": "ClusterIP",
	"sessionAffinity": "None"
  },
  "status": {
	"loadBalancer": {
	}
  }
}
`

var JsonConfigMap = `
{
  "kind": "ConfigMap",
  "apiVersion": "v1",
  "metadata": {
	"name": "otel-collector-config-SINK_ID",
	"creationTimestamp": null
  },
  "data": {
	"config.yaml": "SINK_CONFIG"
  }
}
`

var JsonDeployment = `
{
      "kind": "Deployment",
      "apiVersion": "apps/v1",
      "metadata": {
        "name": "otel-SINK_ID",
        "creationTimestamp": null,
        "labels": {
          "app": "opentelemetry",
          "component": "otel-collector"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "app": "opentelemetry",
            "component": "otel-collector-SINK_ID"
          }
        },
        "template": {
          "metadata": {
            "annotations": {
              "prometheus.io/path": "/metrics",
              "prometheus.io/port": "8888",
              "prometheus.io/scrape": "true"
            },
            "creationTimestamp": null,
            "labels": {
              "app": "opentelemetry",
              "component": "otel-collector-SINK_ID"
            }
          },
          "spec": {
            "volumes": [
              {
                "name": "varlog",
                "hostPath": {
                  "path": "/var/log",
                  "type": ""
                }
              },
              {
                "name": "varlibdockercontainers",
                "hostPath": {
                  "path": "/var/lib/docker/containers",
                  "type": ""
                }
              },
              {
                "name": "data",
                "configMap": {
                  "name": "otel-collector-config-SINK_ID",
                  "defaultMode": 420
                }
              }
            ],
            "containers": [
              {
                "name": "otel-collector",
                "image": "orb-community/otelcol-orb-maestro@latest",
                "ports": [
                  {
                    "containerPort": 13133,
                    "protocol": "TCP"
                  },
                  {
                    "containerPort": 8888,
                    "protocol": "TCP"
                  }
                ],
                "volumeMounts": [
                  {
                    "name": "varlog",
                    "readOnly": true,
                    "mountPath": "/var/log"
                  },
                  {
                    "name": "varlibdockercontainers",
                    "readOnly": true,
                    "mountPath": "/var/lib/docker/containers"
                  },
                  {
                    "name": "data",
                    "readOnly": true,
                    "mountPath": "/etc/otelcol-contrib/config.yaml",
                    "subPath": "config.yaml"
                  }
                ],
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "imagePullPolicy": "IfNotPresent"
              }
            ],
            "restartPolicy": "Always",
            "terminationGracePeriodSeconds": 30,
            "dnsPolicy": "ClusterFirst",
            "securityContext": {
            },
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": "25%",
            "maxSurge": "25%"
          }
        },
        "revisionHistoryLimit": 10,
        "progressDeadlineSeconds": 600
      },
      "status": {
      }
    }
`

func (c *configBuilder) BuildDeploymentConfig(deployment *DeploymentRequest) (string, error) {
	// prepare manifest
	manifest := strings.Replace(k8sOtelCollector, "SINK_ID", deployment.SinkID, -1)
	ctx := context.WithValue(context.Background(), "sink_id", deployment.SinkID)
	config, err := c.ReturnConfigYamlFromSink(ctx, c.kafkaUrl, deployment)
	if err != nil {
		return "", errors.Wrap(errors.New(fmt.Sprintf("failed to build YAML, sink: %s", deployment.SinkID)), err)
	}
	manifest = strings.Replace(manifest, "SINK_CONFIG", config, -1)
	return manifest, nil
}

// ReturnConfigYamlFromSink this is the main method, which will generate the YAML file from the
func (c *configBuilder) ReturnConfigYamlFromSink(_ context.Context, kafkaUrlConfig string, deployment *DeploymentRequest) (string, error) {
	authType := deployment.Config.GetSubMetadata(AuthenticationKey)["type"]
	authTypeStr, ok := authType.(string)
	if !ok {
		return "", errors.New("failed to create config invalid authentication type")
	}
	// TODO move this into somewhere else
	authBuilder := GetAuthService(authTypeStr, c.encryptionService)
	if authBuilder == nil {
		return "", errors.New("invalid authentication type")
	}
	exporterBuilder := FromStrategy(deployment.Backend)
	if exporterBuilder == nil {
		return "", errors.New("invalid backend")
	}
	extensions, extensionName := authBuilder.GetExtensionsFromMetadata(deployment.Config)
	exporters, exporterName := exporterBuilder.GetExportersFromMetadata(deployment.Config, extensionName)
	if exporterName == "" {
		return "", errors.New("failed to build exporter")
	}

	// Add prometheus extension for metrics
	extensions.PProf = &PProfExtension{
		Endpoint: "0.0.0.0:1888",
	}
	serviceConfig := ServiceConfig{
		Extensions: []string{"pprof", extensionName},
		Pipelines: struct {
			Metrics struct {
				Receivers  []string `json:"receivers" yaml:"receivers"`
				Processors []string `json:"processors,omitempty" yaml:"processors,omitempty"`
				Exporters  []string `json:"exporters" yaml:"exporters"`
			} `json:"metrics" yaml:"metrics"`
		}{
			Metrics: struct {
				Receivers  []string `json:"receivers" yaml:"receivers"`
				Processors []string `json:"processors,omitempty" yaml:"processors,omitempty"`
				Exporters  []string `json:"exporters" yaml:"exporters"`
			}{
				Receivers: []string{"kafka"},
				Exporters: []string{exporterName},
			},
		},
	}
	config := OtelConfigFile{
		Receivers: Receivers{
			Kafka: KafkaReceiver{
				Brokers:         []string{kafkaUrlConfig},
				Topic:           fmt.Sprintf("otlp_metrics-%s", deployment.SinkID),
				ProtocolVersion: "2.0.0",
			},
		},
		Extensions: &extensions,
		Exporters:  exporters,
		Service:    serviceConfig,
	}
	marshal, err := yaml.Marshal(&config)
	if err != nil {
		return "", err
	}
	returnedString := "---\n" + string(marshal)
	s := strings.ReplaceAll(returnedString, "\"", "")
	s = strings.ReplaceAll(s, "\n", `\n`)
	return s, nil
}
