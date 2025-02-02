<script setup lang="ts">
import { NConfigProvider, NMessageProvider, NLayout, NLayoutSider, NLayoutContent, NMenu, NIcon } from 'naive-ui'
import { ref, h } from 'vue'
import {
  ServerOutline,
  AppsOutline,
  FolderOutline,
  PersonOutline,
  ListOutline,
  HardwareChipOutline,
  GitNetworkOutline
} from '@vicons/ionicons5'
import Header from './components/Header.vue'
import SystemInfo from './components/SystemInfo.vue'
import SoftwareManager from './components/SoftwareManager.vue'
import FileExplorer from './components/FileExplorer.vue'
import ProcessList from './components/ProcessList.vue'
import PortList from './components/PortList.vue'

const activeKey = ref('system-info')

const renderIcon = (icon: any) => () => h(NIcon, null, { default: () => h(icon) })

const menuOptions = [
  {
    label: '系统信息',
    key: 'system-info',
    icon: renderIcon(HardwareChipOutline)
  },
  {
    label: '软件管理',
    key: 'software',
    icon: renderIcon(AppsOutline)
  },
  {
    label: '文件管理',
    key: 'files',
    icon: renderIcon(FolderOutline)
  },
  {
    label: '进程列表',
    key: 'processes',
    icon: renderIcon(ListOutline)
  },
  {
    label: '端口列表',
    key: 'ports',
    icon: renderIcon(GitNetworkOutline)
  }
]
</script>

<template>
  <n-config-provider>
    <n-message-provider>
      <n-layout position="absolute">
        <Header />
        <n-layout has-sider position="absolute" style="top: 64px">
          <n-layout-sider bordered collapse-mode="width" :collapsed-width="64" :width="240" show-trigger>
            <n-menu v-model:value="activeKey" :options="menuOptions" />
          </n-layout-sider>
          <n-layout-content content-style="padding: 24px;">
            <component :is="activeKey === 'system-info' ? SystemInfo :
              activeKey === 'software' ? SoftwareManager :
                activeKey === 'files' ? FileExplorer :
                  activeKey === 'processes' ? ProcessList :
                    activeKey === 'ports' ? PortList : null" />
          </n-layout-content>
        </n-layout>
      </n-layout>
    </n-message-provider>
  </n-config-provider>
</template>

<style>
body {
  margin: 0;
  padding: 0;
}
</style>
