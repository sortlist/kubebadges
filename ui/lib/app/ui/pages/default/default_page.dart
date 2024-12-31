import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:ui/app/ui/pages/default/badge_card.dart';
import 'package:ui/app/ui/pages/default/card_title.dart';
import 'package:ui/app/controller/badge_controller.dart';
import 'package:ui/app/ui/widgets/badge_dialog.dart';

class DefaultPage extends GetView<BadgeController> {
  const DefaultPage({super.key});

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const CardTitle(title: "Nodes"),
          Obx(
            () => BadgeCard(
              items: [...controller.nodeList],
              kubeBadge: null,
              onTap: (e) {
                showDialog(
                  context: context,
                  builder: (BuildContext context) {
                    return BadgeSettingDialog(kubeBadge: e);
                  },
                );
              },
            ),
          ),
          const CardTitle(title: "Deployments"),
          Obx(() => Row(
            children: [
              const Text("Select Namespace: "),
              const SizedBox(width: 8),
              DropdownButton<String>(
                value: controller.selectedNamespace.value?.name,
                items: controller.namespaceList.keys.map((namespace) {
                  return DropdownMenuItem<String>(
                    value: namespace.name,
                    child: Text(namespace.name),
                  );
                }).toList(),
                onChanged: (val) async {
                  if (val != null) {
                    final selectedBadge = controller.namespaceList.keys
                        .firstWhere((badge) => badge.name == val);
                    controller.selectedNamespace.value = selectedBadge;
                    var deployments = await controller.appService.listDeployments(val, false);
                    if (!deployments.status.hasError && deployments.body!.isNotEmpty) {
                      controller.namespaceList[selectedBadge] = deployments.body!;
                      controller.refreshNamespaceList();
                    }
                  }
                },
                hint: const Text("Choose a namespace"),
              ),
            ],
          )),
          Obx(() => controller.selectedNamespace.value != null
              ? BadgeCard(
                  items: controller.namespaceList[controller.selectedNamespace.value] ?? [],
                  kubeBadge: controller.selectedNamespace.value,
                  onTap: (e) {
                    showDialog(
                      context: context,
                      builder: (BuildContext context) {
                        return BadgeSettingDialog(kubeBadge: e);
                      },
                    );
                  },
                )
              : const SizedBox()),
          const CardTitle(title: "Kustomizations"),
          Obx(() => Row(
            children: [
              const Text("Select Namespace: "),
              const SizedBox(width: 8),
              DropdownButton<String>(
                value: controller.selectedKustomizationNamespace.value?.name,
                items: controller.kustomizationNamespaceList.keys.map((namespace) {
                  return DropdownMenuItem<String>(
                    value: namespace.name,
                    child: Text(namespace.name),
                  );
                }).toList(),
                onChanged: (val) async {
                  if (val != null) {
                    final selectedBadge = controller.kustomizationNamespaceList.keys
                        .firstWhere((badge) => badge.name == val);
                    controller.selectedKustomizationNamespace.value = selectedBadge;
                    var kustomizations = await controller.appService.listKustomizations(val, false);
                    if (!kustomizations.status.hasError && kustomizations.body!.isNotEmpty) {
                      controller.kustomizationNamespaceList[selectedBadge] = kustomizations.body!;
                      controller.refreshKustomizationNamespaceList();
                    }
                  }
                },
                hint: const Text("Choose a namespace"),
              ),
            ],
          )),
          Obx(() => controller.selectedKustomizationNamespace.value != null
              ? BadgeCard(
                  items: controller.kustomizationNamespaceList[controller.selectedKustomizationNamespace.value] ?? [],
                  kubeBadge: controller.selectedKustomizationNamespace.value,
                  onTap: (e) {
                    showDialog(
                      context: context,
                      builder: (BuildContext context) {
                        return BadgeSettingDialog(kubeBadge: e);
                      },
                    );
                  },
                )
              : const SizedBox()),
          const CardTitle(title: "PostgreSQL"),
          Obx(() => Row(
            children: [
              const Text("Select Namespace: "),
              const SizedBox(width: 8),
              DropdownButton<String>(
                value: controller.selectedPostgresqlNamespace.value?.name,
                items: controller.postgresqlNamespaceList.keys.map((namespace) {
                  return DropdownMenuItem<String>(
                    value: namespace.name,
                    child: Text(namespace.name),
                  );
                }).toList(),
                onChanged: (val) async {
                  if (val != null) {
                    final selectedBadge = controller.postgresqlNamespaceList.keys
                        .firstWhere((badge) => badge.name == val);
                    controller.selectedPostgresqlNamespace.value = selectedBadge;
                    var postgresqls = await controller.appService.listPostgresqls(val, false);
                    if (!postgresqls.status.hasError && postgresqls.body!.isNotEmpty) {
                      controller.postgresqlNamespaceList[selectedBadge] = postgresqls.body!;
                      controller.refreshPostgresqlNamespaceList();
                    }
                  }
                },
                hint: const Text("Choose a namespace"),
              ),
            ],
          )),
          Obx(() => controller.selectedPostgresqlNamespace.value != null
              ? BadgeCard(
                  items: controller.postgresqlNamespaceList[controller.selectedPostgresqlNamespace.value] ?? [],
                  kubeBadge: controller.selectedPostgresqlNamespace.value,
                  onTap: (e) {
                    showDialog(
                      context: context,
                      builder: (BuildContext context) {
                        return BadgeSettingDialog(kubeBadge: e);
                      },
                    );
                  },
                )
              : const SizedBox()),
          const CardTitle(title: "Jobs"),
          Obx(() => Row(
            children: [
              const Text("Select Namespace: "),
              const SizedBox(width: 8),
              DropdownButton<String>(
                value: controller.selectedJobNamespace.value?.name,
                items: controller.jobNamespaceList.keys.map((namespace) {
                  return DropdownMenuItem<String>(
                    value: namespace.name,
                    child: Text(namespace.name),
                  );
                }).toList(),
                onChanged: (val) async {
                  if (val != null) {
                    final selectedBadge = controller.jobNamespaceList.keys
                        .firstWhere((badge) => badge.name == val);
                    controller.selectedJobNamespace.value = selectedBadge;
                    var jobs = await controller.appService.listJobs(val, false);
                    if (!jobs.status.hasError && jobs.body!.isNotEmpty) {
                      controller.jobNamespaceList[selectedBadge] = jobs.body!;
                      controller.refreshJobNamespaceList();
                    }
                  }
                },
                hint: const Text("Choose a namespace"),
              ),
            ],
          )),
          Obx(() => controller.selectedJobNamespace.value != null
              ? BadgeCard(
                  items: controller.jobNamespaceList[controller.selectedJobNamespace.value] ?? [],
                  kubeBadge: controller.selectedJobNamespace.value,
                  onTap: (e) {
                    showDialog(
                      context: context,
                      builder: (BuildContext context) {
                        return BadgeSettingDialog(kubeBadge: e);
                      },
                    );
                  },
                )
              : const SizedBox()),
        ],
      ),
    );
  }
}
