# Java Http Server with OpenTelemetry instrumentation

This app listens on port `3000` and exposes a single endpoint at `/` that responds with the string "hello from java\n
this is the end of the journey... for today".

The OpenTelemetry Java HTTP Example was used as a reference for this sample.
https://github.com/open-telemetry/opentelemetry-java/tree/master/examples/http


## Running the application
TODO

## Add the relevant dependencies and repositories to pom.xml

```diff
  <!-- library dependencies -->
+  <dependencies>
+    <dependency>
+      <groupId>io.opentelemetry</groupId>
+      <artifactId>opentelemetry-sdk</artifactId>
+      <version>0.4.0-SNAPSHOT</version>
+    </dependency>
+    <dependency>
+      <groupId>io.opentelemetry</groupId>
+      <artifactId>opentelemetry-exporters-logging</artifactId>
+      <version>0.4.0-SNAPSHOT</version>
+    </dependency>
+  </dependencies>
+  <repositories>
+    <repository>
+      <id>oss.sonatype.org-snapshot</id>
+      <url>https://oss.jfrog.org/artifactory/oss-snapshot-local</url>
+    </repository>
+  </repositories>
```

## Import the packages required for instrumenting your Java app

```diff
+import io.opentelemetry.OpenTelemetry;
+import io.opentelemetry.common.AttributeValue;
+import io.opentelemetry.context.ContextUtils;
+import io.opentelemetry.context.Scope;
+import io.opentelemetry.context.propagation.HttpTextFormat;
+import io.opentelemetry.exporters.logging.LoggingSpanExporter;
+import io.opentelemetry.sdk.OpenTelemetrySdk;
+import io.opentelemetry.sdk.trace.TracerSdkProvider;
+import io.opentelemetry.sdk.trace.export.SimpleSpansProcessor;
+import io.opentelemetry.trace.*;

```

## Initiate the tracer, the logging exporter and invoke it during the Main class initializer.

```diff

+  // OTel API
+  private static Tracer tracer =
+      OpenTelemetry.getTracerProvider().get("Main");
+  // Export traces to log.
+  private static LoggingSpanExporter loggingExporter = new LoggingSpanExporter();

+  private void initTracer() {
+    // Get the tracer
+    TracerSdkProvider tracerProvider = OpenTelemetrySdk.getTracerProvider();
+    // Show that multiple exporters can be used

+    // Set to export the traces also to a log file.
+    tracerProvider.addSpanProcessor(SimpleSpansProcessor.newBuilder(loggingExporter).build());
+  }
```

```diff
private Main(int port) throws IOException {
+    initTracer();

```

## Instrument the HTTP Handler

```diff
    @Override
    public void handle(HttpExchange he) throws IOException {
+      Span span = tracer.spanBuilder("java app").setSpanKind(Span.Kind.SERVER);
+      span.setAttribute("http.url", "https://important-dusty-apparatus.glitch.me/");

       // Process the request
      String response = "hello from java\n this is the end of the journey... for today";
      he.sendResponseHeaders(200, response.length());
      OutputStream os = he.getResponseBody();
      os.write(response.getBytes(Charset.defaultCharset()));
      os.close();
      System.out.println("Served Client: " + he.getRemoteAddress());

+      span.setAttribute("response", response);

+      // Everything works fine in this example
+      span.setStatus(Status.OK);

+      // Close the span
+      span.end();
    }
```