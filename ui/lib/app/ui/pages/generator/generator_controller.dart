import 'package:get/get.dart';
import 'package:ui/app/model/model.dart';
import 'package:ui/app/service/app_service.dart';

// Simple enum for supported types
enum ResourceType {
  node,
  deployment,
  kustomization,
  // more can be added
}

class GeneratorController extends GetxController {
  final AppService appService = Get.find();

  var resourceType = ResourceType.node.obs;
  var namespace = ''.obs;
  var resourceList = <KubeBadge>[].obs;
  var selectedItem = Rxn<KubeBadge>();

  // For the list of namespaces
  var namespaces = <KubeBadge>[].obs;

  @override
  void onInit() {
    super.onInit();
    loadNamespaces();
    loadResources(); // node by default
  }

  void loadNamespaces() async {
    final resp = await appService.listNamespace(false);
    if (!resp.status.hasError) {
      namespaces.assignAll(resp.body!);
    }
  }

  // Load the list according to type
  void loadResources() async {
    resourceList.clear();
    selectedItem.value = null;
    if (resourceType.value == ResourceType.node) {
      final resp = await appService.listNodes(false);
      if (!resp.status.hasError) {
        resourceList.assignAll(resp.body!);
      }
    } else if (resourceType.value == ResourceType.deployment) {
      // namespace required
      if (namespace.value.isEmpty) return;
      final resp = await appService.listDeployments(namespace.value, false);
      if (!resp.status.hasError) {
        resourceList.assignAll(resp.body!);
      }
    } else if (resourceType.value == ResourceType.kustomization) {
      if (namespace.value.isEmpty) return;
      final resp = await appService.get('/api/kustomizations/${namespace.value}', decoder: (data) {
        return (data as List).map((e) => KubeBadge.fromJson(e)).toList();
      });
      if (!resp.hasError) {
        resourceList.assignAll(resp.body!);
      }
    }
  }

  String getBadgeURL() {
    if (selectedItem.value == null) return '';
    return appService.getBadgeBaseURL() + selectedItem.value!.badge;
  }
}
