package com.fgsantana.ms_java_prometheus;

import io.micrometer.core.instrument.binder.jvm.ClassLoaderMetrics;
import io.micrometer.core.instrument.binder.jvm.JvmGcMetrics;
import io.micrometer.core.instrument.binder.jvm.JvmMemoryMetrics;
import io.micrometer.core.instrument.binder.jvm.JvmThreadMetrics;
import io.micrometer.core.instrument.binder.system.ProcessorMetrics;
import io.micrometer.core.instrument.binder.system.UptimeMetrics;
import io.micrometer.prometheus.PrometheusConfig;
import io.micrometer.prometheus.PrometheusMeterRegistry;
import io.vertx.core.*;
import io.vertx.core.http.HttpMethod;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.Router;


import java.util.HashMap;



public class Server extends AbstractVerticle {

  @Override
  public void start(Promise<Void> startPromise) throws Exception {
    Router router = Router.router(vertx);

    PrometheusMeterRegistry registry = new PrometheusMeterRegistry(PrometheusConfig.DEFAULT);

    registry.config().commonTags("application","ms-kt-prom");

    new ClassLoaderMetrics().bindTo(registry);
    new JvmMemoryMetrics().bindTo(registry);
    new JvmGcMetrics().bindTo(registry);
    new ProcessorMetrics().bindTo(registry);
    new JvmThreadMetrics().bindTo(registry);
    new UptimeMetrics().bindTo(registry);


    router.route(HttpMethod.GET,"/metrics").handler( ctx -> ctx.response()
      .putHeader("Content-Type","text/plain")
      .end(registry.scrape()));

    router.route(HttpMethod.GET,"/test").handler( ctx -> {
      ctx.response().putHeader("content-type","application/json");
      ctx.json(new JsonObject().put("method","GET"));
    });

    router.route(HttpMethod.POST,"/test").handler( ctx -> {
      ctx.response().putHeader("content-type","application/json");
      HashMap<String,String> map = new HashMap<>();
      map.put("method","POST");
      ctx.json(map);
    });

    vertx.createHttpServer().requestHandler(router).listen(8084).onSuccess( h -> System.out.println("Listening in 8084"));

  }
}
