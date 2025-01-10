<script lang="ts">
	import { on } from "svelte/events";

    let { passwordItem = $bindable() } = $props()

    function getComponent(type: string) {
        switch(type) {
            case 'password':
                return 
        }
    }
    $effect(() => {
        passwordItem
    })
</script>
<div class=" bg-page_faint w-full border-t-2 border-border_primary content-center">
    
</div>
