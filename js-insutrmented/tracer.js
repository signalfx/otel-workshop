'use strict';

const opentelemetry = require('@opentelemetry/api');
const { NodeTracerProvider } = require('@opentelemetry/node');
const { SimpleSpanProcessor } = require('@opentelemetry/tracing');
const { CollectorExporter } = require('@opentelemetry/exporter-collector');

const EXPORTER = process.env.EXPORTER || '';


function initTracer(serviceName) {
  const provider = new NodeTracerProvider()
  const exporter = new CollectorExporter({serviceName});
  provider.addSpanProcessor(new SimpleSpanProcessor(exporter));

  // Initialize the OpenTelemetry APIs to use the NodeTracerProvider bindings
  provider.register();

  opentelemetry.trace.getTracer('express-example');
};
