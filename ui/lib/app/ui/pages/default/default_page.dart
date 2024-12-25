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
              DropdownButton<KubeBadge>(
                value: controller.selectedNamespace.value,
                items: controller.namespaceList.keys.map((namespace) {
                  return DropdownMenuItem<KubeBadge>(
                    value: namespace,
                    child: Text(namespace.name),
                  );
                }).toList(),
                onChanged: (val) async {
                  if (val != null) {
                    controller.selectedNamespace.value = val;
                    var deployments = await controller.appService.listDeployments(val.name, false);
                    if (!deployments.status.hasError && deployments.body!.isNotEmpty) {
                      controller.namespaceList[val] = deployments.body!;
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
        ],
      ),
    );
  }
}
