<script lang="ts">
    var {className, children} = $props()
</script>

<div class="sm:flex sm:justify-between sm:gap-4 {className}">
    <div>
        <h3 class="text-gray-900 text-lg font-bold sm:text-xl">{@render children?.()}</h3>
    </div>
</div>