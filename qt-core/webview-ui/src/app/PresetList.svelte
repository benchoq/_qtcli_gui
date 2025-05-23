<!--
Copyright (C) 2025 The Qt Company Ltd.
SPDX-License-Identifier: LicenseRef-Qt-Commercial OR LGPL-3.0-only 
-->

<script lang="ts">
  import { Listgroup, ListgroupItem, P } from "flowbite-svelte";

  import { presets } from "./states.svelte";
  import * as viewlogic from "./viewlogic.svelte";

  const adjustSelectedIndex = (offset: number) => {
    if (!presets.selected || presets.selectedIndex < 0) {
      return;
    }

    let candidate = presets.selectedIndex + offset;
    candidate = Math.max(0, candidate);
    candidate = Math.min(candidate, presets.all.length - 1);

    if (candidate != presets.selectedIndex) {
      viewlogic.setSelectedPreset(presets.all[candidate], candidate);
    }
  };

  const onKeyPressed = (e: KeyboardEvent) => {
    if (e.key === "ArrowUp") {
      adjustSelectedIndex(-1);
    } else if (e.key === "ArrowDown") {
      adjustSelectedIndex(+1);
    } else if (e.key === "Enter") {
      viewlogic.moveWizardPage(1)
    } else {
      return;
    }

    const el = e.currentTarget as HTMLElement;
    if (el) {
      const items = el?.querySelectorAll("button");
      const item = items[presets.selectedIndex];
      if (item instanceof HTMLButtonElement) {
        item.focus();
      }
    }
  };
</script>

<div class="flex flex-col">
  <Listgroup
    active
    class="flex-grow overflow-y-auto qt-list"
    onkeydown={onKeyPressed}
    tabindex={0}
  >
    {#each presets.all as preset, index}
      <ListgroupItem
        class="qt-list-item flex flex-row"
        currentClass="selected"
        current={presets.selectedIndex === index}
        on:click={() => {
          viewlogic.setSelectedPreset(preset, index);
        }}
      >
        <div class="flex-1"
          role="listitem"
          on:dblclick={() => { viewlogic.moveWizardPage(1); }}
        >
          {viewlogic.createPresetDisplayText(preset)}
        </div>

        {#if !preset.name.startsWith("@")}
          <P size="xs" class="qt-label">{preset.meta.title}</P>
        {:else}
          <div></div>
        {/if}
      </ListgroupItem>
    {/each}
  </Listgroup>
</div>
