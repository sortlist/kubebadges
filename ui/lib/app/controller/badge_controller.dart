import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:get/get.dart';
import 'package:ui/app/model/model.dart';
import 'package:ui/app/service/app_service.dart';

class BadgeController extends GetxController {
  AppService appService = Get.find();

  final _nodeList = <KubeBadge>[].obs;
  List<KubeBadge> get nodeList => _nodeList;
  set nodeList(List<KubeBadge> value) => _nodeList.value = value;

  final _namespaceList = <KubeBadge, List<KubeBadge>>{}.obs;
  Map<KubeBadge, List<KubeBadge>> get namespaceList => _namespaceList;
  set namespaceList(Map<KubeBadge, List<KubeBadge>> value) =>
      _namespaceList.value = value;

  final selectedNamespace = Rxn<KubeBadge>();

  final _kustomizationNamespaceList = <KubeBadge, List<KubeBadge>>{}.obs;
  Map<KubeBadge, List<KubeBadge>> get kustomizationNamespaceList => _kustomizationNamespaceList;
  set kustomizationNamespaceList(Map<KubeBadge, List<KubeBadge>> value) =>
      _kustomizationNamespaceList.value = value;

  final selectedKustomizationNamespace = Rxn<KubeBadge>();

  final _postgresqlNamespaceList = <KubeBadge, List<KubeBadge>>{}.obs;
  Map<KubeBadge, List<KubeBadge>> get postgresqlNamespaceList => _postgresqlNamespaceList;
  set postgresqlNamespaceList(Map<KubeBadge, List<KubeBadge>> value) =>
      _postgresqlNamespaceList.value = value;

  final selectedPostgresqlNamespace = Rxn<KubeBadge>();

  final _jobNamespaceList = <KubeBadge, List<KubeBadge>>{}.obs;
  Map<KubeBadge, List<KubeBadge>> get jobNamespaceList => _jobNamespaceList;
  set jobNamespaceList(Map<KubeBadge, List<KubeBadge>> value) =>
      _jobNamespaceList.value = value;

  final selectedJobNamespace = Rxn<KubeBadge>();

  void refreshNamespaceList() {
    _namespaceList.refresh();
  }

  void refreshKustomizationNamespaceList() {
    _kustomizationNamespaceList.refresh();
  }

  void refreshPostgresqlNamespaceList() {
    _postgresqlNamespaceList.refresh();
  }

  void refreshJobNamespaceList() {
    _jobNamespaceList.refresh();
  }

  BadgeController() {
    loadData(false);
  }

  void loadData(bool force) async {
    listNodes(force);
    listNamespace(force);
  }

  void listNodes(bool force) async {
    namespaceList.clear();
    var response = await appService.listNodes(force);
    if (!response.status.hasError) {
      nodeList = response.body!;
    }
  }

  void listNamespace(bool force) async {
    EasyLoading.show(status: 'Loading...');
    namespaceList.clear();
    kustomizationNamespaceList.clear();
    postgresqlNamespaceList.clear();
    jobNamespaceList.clear();
    var namespaces = await appService.listNamespace(force);
    if (!namespaces.status.hasError) {
      for (var namespace in namespaces.body!) {
        namespaceList[namespace] = [];
        kustomizationNamespaceList[namespace] = [];
        postgresqlNamespaceList[namespace] = [];
        jobNamespaceList[namespace] = [];
      }
    }
    EasyLoading.dismiss();
  }

  Future<bool> updateBadgePublic(KubeBadge kubeBadge, bool allowed) async {
    EasyLoading.show(status: 'Loading...');
    var response = await appService.updateBadge({
      "key": kubeBadge.key,
      "allowed": allowed,
    });
    if (response.status.hasError) {
      final Map<String, dynamic> error = response.body;
      EasyLoading.dismiss();
      EasyLoading.showError(
        'Badge update failed : ${error['error']}',
        duration: const Duration(milliseconds: 2500),
      );
      return false;
    }

    int index = nodeList.indexWhere((element) => element.key == kubeBadge.key);
    if (index != -1) {
      nodeList[index] = kubeBadge.copyWith(allowed: allowed);
    }

    for (var namespace in namespaceList.keys) {
      int index = namespaceList[namespace]!
          .indexWhere((element) => element.key == kubeBadge.key);
      if (index != -1) {
        var updatedList = List<KubeBadge>.from(namespaceList[namespace]!);
        updatedList[index] = kubeBadge.copyWith(allowed: allowed);
        namespaceList[namespace] = updatedList;
      }
    }
    namespaceList = Map.from(namespaceList);
    EasyLoading.dismiss();
    EasyLoading.showToast(
      'Badge updated success',
      duration: const Duration(seconds: 1),
      toastPosition: EasyLoadingToastPosition.bottom,
    );
    return true;
  }

  String getBadgeBaseURL() {
    return appService.getBadgeBaseURL();
  }
}
