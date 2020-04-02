import requests
from flask import Flask
from opentelemetry import trace
from opentelemetry.ext.otcollector.trace_exporter import CollectorSpanExporter
from opentelemetry.ext import http_requests
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchExportSpanProcessor
from opentelemetry.ext.wsgi import OpenTelemetryMiddleware


# configure tracer and exporter
ot_exporter = CollectorSpanExporter(
    service_name="py-service",
    endpoint="localhost:55678",
)

provider = TracerProvider()
trace.set_tracer_provider(provider)
tracer = trace.get_tracer(__name__)
provider.add_span_processor(BatchExportSpanProcessor(ot_exporter))

# instrument http client
http_requests.enable(provider)

# create and instrument flask server
app = Flask(__name__)
app.wsgi_app = OpenTelemetryMiddleware(app.wsgi_app)


@app.route("/")
def hello():
    response = fetch_from_java()
    return "hello from python\n" + response


def fetch_from_java():
    r = requests.get('http://localhost:8081/')
    return r.text


if __name__ == "__main__":
    app.run(debug=True, port=8082)
