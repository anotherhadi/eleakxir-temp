<script lang="ts" module>
  interface NavItem {
    title: string;
    url: string;
    items?: NavItem[];
    isActive?: boolean;
    icon?: any;
  }

  interface Nav {
    navMain: NavItem[];
  }

  const data: Nav = {
    navMain: [
      {
        title: "Home",
        url: "/",
        icon: Home,
      },
      {
        title: "Search",
        url: "/search",
        icon: Search,
      },
      {
        title: "Data wells",
        url: "/dataleaks",
        icon: Database,
      },
      {
        title: "Parquet files",
        url: "/parquet",
        icon: FileSearch,
      },
    ],
  };
</script>

<script lang="ts">
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import { Database, FileSearch, Github, Home, Search } from "@lucide/svelte";
  import { isActiveLink } from "sv-router";
  import type { ComponentProps } from "svelte";
  import ModeToggle from "../mode-toggle.svelte";
  import Button from "../ui/button/button.svelte";
  let {
    ref = $bindable(null),
    ...restProps
  }: ComponentProps<typeof Sidebar.Root> = $props();
</script>

<Sidebar.Root variant="floating" {...restProps}>
  <Sidebar.Header>
    <Sidebar.Menu>
      <Sidebar.MenuItem>
        <Sidebar.MenuButton size="lg">
          {#snippet child({ props })}
            <a href="/" {...props}>
              <div
                class="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg"
              >
                <img src="/favicon.svg" alt="logo" class="w-4/5" />
              </div>
              <div class="flex flex-col gap-0.5 leading-none">
                <span class="font-medium capitalize">Eleakxir</span>
              </div>
            </a>
          {/snippet}
        </Sidebar.MenuButton>
      </Sidebar.MenuItem>
    </Sidebar.Menu>
  </Sidebar.Header>
  <Sidebar.Content>
    <Sidebar.Group>
      <Sidebar.Menu class="gap-2">
        {#each data.navMain as item (item.title)}
          <Sidebar.MenuItem>
            <Sidebar.MenuButton>
              {#snippet child({ props })}
                <a
                  href={item.url}
                  class="font-medium flex gap-2 items-center"
                  use:isActiveLink={{ className: "bg-accent" }}
                  {...props}
                >
                  <item.icon />
                  {item.title}
                </a>
              {/snippet}
            </Sidebar.MenuButton>
            {#if item.items?.length}
              <Sidebar.MenuSub class="ml-0 border-l-0 px-1.5">
                {#each item.items as subItem (subItem.title)}
                  <Sidebar.MenuSubItem>
                    <Sidebar.MenuSubButton isActive={subItem.isActive}>
                      {#snippet child({ props })}
                        <a href={subItem.url} {...props}>
                          {subItem.title}
                          {subItem.title}</a
                        >
                      {/snippet}
                    </Sidebar.MenuSubButton>
                  </Sidebar.MenuSubItem>
                {/each}
              </Sidebar.MenuSub>
            {/if}
          </Sidebar.MenuItem>
        {/each}
      </Sidebar.Menu>
    </Sidebar.Group>
  </Sidebar.Content>
  <Sidebar.Footer>
    <div class="flex flex-wrap gap-2 items-center justify-center">
      <Button
        variant="ghost"
        size="icon"
        href="https://github.com/anotherhadi/eleakxir"><Github /></Button
      >
      <ModeToggle />
    </div>
  </Sidebar.Footer>
</Sidebar.Root>
