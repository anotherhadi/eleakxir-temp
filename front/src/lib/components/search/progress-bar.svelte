<script lang="ts">
  import Progress from "../ui/progress/progress.svelte";

  const { progress } = $props();
</script>

<div class="my-10 flex flex-col gap-5 items-center">
  <p class="text-muted-foreground animate-pulse">
    Loading... {Math.round(progress)}%</p>
  <Progress
    value={progress}
    class="w-full max-w-lg transition-transform"
  />
</div>
