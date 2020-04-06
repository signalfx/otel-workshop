import io.grpc.Context;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;
import io.opentelemetry.OpenTelemetry;
import io.opentelemetry.common.AttributeValue;
import io.opentelemetry.context.ContextUtils;
import io.opentelemetry.context.Scope;
import io.opentelemetry.context.propagation.HttpTextFormat;
import io.opentelemetry.exporters.logging.LoggingSpanExporter;
import io.opentelemetry.sdk.OpenTelemetrySdk;
import io.opentelemetry.sdk.trace.TracerSdkProvider;
import io.opentelemetry.sdk.trace.export.SimpleSpansProcessor;
import io.opentelemetry.trace.*;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.Charset;

public class MainInstrumented {

    private class HelloHandler implements HttpHandler {

        @Override
        public void handle(HttpExchange he) throws IOException {
            Span span = tracer.spanBuilder("java app").setSpanKind(Span.Kind.SERVER).startSpan();
            span.setAttribute("http.url", "https://important-dusty-apparatus.glitch.me/");

            // Process the request
            String response = "hello from java\n this is the end of the journey... for today";
            he.sendResponseHeaders(200, response.length());
            OutputStream os = he.getResponseBody();
            os.write(response.getBytes(Charset.defaultCharset()));
            os.close();
            System.out.println("Served Client: " + he.getRemoteAddress());

            span.setAttribute("response", response);

            // Everything works fine in this example
            span.setStatus(Status.OK);

            // Close the span
            span.end();
        }

    }

    private com.sun.net.httpserver.HttpServer server;
    private static int port = 3000;
    private MainInstrumented() throws IOException {
        this(port);
    }

    private MainInstrumented(int port) throws IOException {
        initTracer();
        server = com.sun.net.httpserver.HttpServer.create(new InetSocketAddress(port), 0);
        // Test urls
        server.createContext("/", new HelloHandler());
        server.start();
        System.out.println("Server ready on http://127.0.0.1:" + port);
    }

    // OTel API
    private static Tracer tracer =
            OpenTelemetry.getTracerProvider().get("Main");
    // Export traces to log
    private static LoggingSpanExporter loggingExporter = new LoggingSpanExporter();

    private void initTracer() {
        // Get an instace of tracer
        TracerSdkProvider tracerProvider = OpenTelemetrySdk.getTracerProvider();

        // Set to export the traces also to a log file.
        // Currently, Java doesn't export in the OpenTelemetry protocol format. It is expected to land soon.
        tracerProvider.addSpanProcessor(SimpleSpansProcessor.newBuilder(loggingExporter).build());
    }

    private void stop() {
        server.stop(0);
    }

    /**
     * Main method to run the example.
     *
     * @param args It is not required.
     * @throws Exception Something might go wrong.
     */
    public static void main(String[] args) throws Exception {
        final MainInstrumented s = new MainInstrumented();
        // Gracefully close the server
        Runtime.getRuntime()
                .addShutdownHook(
                        new Thread() {
                            @Override
                            public void run() {
                                s.stop();
                            }
                        });
    }
}
