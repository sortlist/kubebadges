import 'dart:convert';

import 'package:get/get.dart';
import 'package:ui/app/config/constant.dart';
import 'package:ui/app/model/model.dart';

class Api extends GetConnect {
  @override
  void onInit() {
    super.onInit();
    httpClient.baseUrl = Constants.baseAPI;
  }

  Future<Response<List<KubeBadge>>> listNodes(bool force) {
    return get('/api/nodes?force=$force', decoder: (data) {
      return (data as List).map((item) => KubeBadge.fromJson(item)).toList();
    });
  }

  Future<Response<List<KubeBadge>>> listNamespace(bool force) {
    return get('/api/namespaces?force=$force', decoder: (data) {
      return (data as List).map((item) => KubeBadge.fromJson(item)).toList();
    });
  }

  Future<Response<List<KubeBadge>>> listDeployments(String name, bool force) {
    return get('/api/deployments/$name?force=$force', decoder: (data) {
      return (data as List).map((item) => KubeBadge.fromJson(item)).toList();
    });
  }

  Future<Response> updateBadge(Map<String, dynamic> data) {
    return post('/api/badge', data);
  }

  Future<Response<KubeBadgeConfig>> getBadgeConfig() {
    return get('/api/config', decoder: (data) {
      return KubeBadgeConfig.fromJson(data);
    });
  }

  Future<Response<KubeBadgeConfig>> updateBadgeConfig(KubeBadgeConfig data) {
    return post('/api/config', jsonEncode(data), decoder: (data) {
      return KubeBadgeConfig.fromJson(data);
    });
  }

  Future<Response<List<KubeBadge>>> listKustomizations(String namespace, bool force) {
    return get('/api/kustomizations/$namespace?force=$force', decoder: (data) {
      return (data as List).map((item) => KubeBadge.fromJson(item)).toList();
    });
  }

  Future<Response<List<KubeBadge>>> listPostgresqls(String namespace, bool force) {
    return get('/api/postgresqls/$namespace?force=$force', decoder: (data) {
      return (data as List).map((item) => KubeBadge.fromJson(item)).toList();
    });
  }

  Future<Response<List<KubeBadge>>> listJobs(String namespace, bool force) {
    return get('/api/jobs/$namespace?force=$force', decoder: (data) {
      return (data as List).map((item) => KubeBadge.fromJson(item)).toList();
    });
  }
}
