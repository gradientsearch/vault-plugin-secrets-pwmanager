<script lang="ts">
    let {selectedVault = $bindable()} = $props()
</script>

<div class="w-full max-w-96 bg-page_faint border-2 border-border_primary overflow-y-scroll">
<header class="sticky top-0 border-b-2 border-border_primary p-2">
    
    <h1 class="text-base">{selectedVault?.Name} Vault</h1>
    
</header>

</div>
