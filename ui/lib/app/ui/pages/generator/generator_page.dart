import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:ui/app/ui/pages/generator/generator_controller.dart';
import 'package:flutter_easyloading/flutter_easyloading.dart';
import 'package:clipboard/clipboard.dart';
import 'package:ui/app/ui/widgets/badge.dart';

class GeneratorPage extends GetView<GeneratorController> {
  const GeneratorPage({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    // Wrap in SingleChildScrollView if needed
    return SingleChildScrollView(
      child: Obx(
        () => Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Badges Generator',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 16),

            // Choix du ResourceType
            Row(
              children: [
                const Text("Resource Type: "),
                const SizedBox(width: 8),
                DropdownButton<ResourceType>(
                  value: controller.resourceType.value,
                  items: ResourceType.values.map((rt) {
                    return DropdownMenuItem<ResourceType>(
                      value: rt,
                      child: Text(rt.name),
                    );
                  }).toList(),
                  onChanged: (val) {
                    if (val != null) {
                      controller.resourceType.value = val;
                      // reload when resource type changes
                      controller.loadResources();
                    }
                  },
                ),
              ],
            ),
            const SizedBox(height: 16),

            // Si ResourceType == deployment, kustomization ou postgresql, on a besoin d'un namespace
            if (controller.resourceType.value == ResourceType.deployment ||
                controller.resourceType.value == ResourceType.kustomization ||
                controller.resourceType.value == ResourceType.postgresql) ...[
              const Text("Select Namespace:"),
              const SizedBox(height: 8),
              DropdownButton<String>(
                value: controller.namespace.value.isEmpty ? null : controller.namespace.value,
                items: controller.namespaces.map((ns) {
                  return DropdownMenuItem<String>(
                    value: ns.name,
                    child: Text(ns.name),
                  );
                }).toList(),
                onChanged: (val) {
                  if (val != null) {
                    controller.namespace.value = val;
                    controller.loadResources();
                  }
                },
                hint: const Text("Pick a Namespace"),
              ),
              const SizedBox(height: 16),
            ],

            // Liste des ressources
            const Text("Select Resource:"),
            const SizedBox(height: 8),
            DropdownButton<String>(
              value: controller.selectedItem.value?.name,
              items: controller.resourceList.map((k) {
                return DropdownMenuItem<String>(
                  value: k.name,
                  child: Text(k.name),
                );
              }).toList(),
              onChanged: (val) {
                if (val != null) {
                  // find the corresponding KubeBadge
                  final found = controller.resourceList.firstWhereOrNull((x) => x.name == val);
                  controller.selectedItem.value = found;
                }
              },
              hint: const Text("Pick a resource"),
            ),

            const SizedBox(height: 16),

            // Affiche l'URL
            const Text(
              "Badge URL:",
              style: TextStyle(fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 8),
            SelectableText(controller.getBadgeURL()),
            const SizedBox(height: 16),
            if (controller.selectedItem.value != null) ...[
              const Text(
                "Preview:",
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 8),
              KubeBadgeView(
                badge: controller.selectedItem.value!,
                onTap: () {}, // No action needed on tap for preview
              ),
              const SizedBox(height: 16),
            ],
            ElevatedButton(
              onPressed: () {
                final text = controller.getBadgeURL();
                if (text.isNotEmpty) {
                  FlutterClipboard.copy(text).then((value) {
                    EasyLoading.showToast('Copied to clipboard!');
                  });
                }
              },
              child: const Text("Copy URL"),
            )
          ],
        ),
      ),
    );
  }
}
