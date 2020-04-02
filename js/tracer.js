'use strict';

const opentelemetry = require('@opentelemetry/api');
const { NodeTracerProvider } = require('@opentelemetry/node');
const { SimpleSpanProcessor } = require('@opentelemetry/tracing');
// const { ZipkinExporter } = require('@opentelemetry/exporter-zipkin');
const { CollectorExporter } = require('@opentelemetry/exporter-collector');

const EXPORTER = process.env.EXPORTER || '';


module.exports = (serviceName) => {
  const provider = new NodeTracerProvider()
      /*
  const provider = new NodeTracerProvider({
    plugins: {
      express: {
        enabled: true,
        path: '@opentelemetry/plugin-express',
      },
      http: {
        enabled: true,
        path: '@opentelemetry/plugin-http',
      },
    },
  });
    */

  // const exporter = new ZipkinExporter({serviceName});
  const exporter = new CollectorExporter({serviceName});

  provider.addSpanProcessor(new SimpleSpanProcessor(exporter));

  // Initialize the OpenTelemetry APIs to use the NodeTracerProvider bindings
  provider.register();

  return opentelemetry.trace.getTracer('express-example');
};
