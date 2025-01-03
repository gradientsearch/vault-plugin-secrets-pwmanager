<script lang="ts">
	import { base } from "$app/paths";

    base
</script>

<header class="p-4 bg-page_primary border-b-border_primary border-b-2">
	<h1 class="text-start font-bold">Create Password</h1>
</header>

<div class="grid grid-cols-3 gap-4 p-4">
    <button class="relative min-h-32  bg-page_faint shadow-lg flex md:flex-row flex-col justify-items-start items-center hover:bg-surface_interactive_hover">
        <img class="h-16" src='{base}/icons/key.svg' alt="key icon">
        <span class="md:absolute  left-0 right-0  flex-1 text-center">Login</span>
    </button>

    <button class="relative min-h-32  bg-page_faint shadow-lg flex md:flex-row flex-col justify-items-start items-center hover:bg-surface_interactive_hover">
        <img class="h-16" src='{base}/icons/note.png' alt="note icon">
        <span class="md:absolute  left-0 right-0  flex-1 text-center">Note</span>
    </button>

    <button class="relative min-h-32  bg-page_faint shadow-lg flex md:flex-row flex-col justify-items-start items-center hover:bg-surface_interactive_hover">
        <img class="h-16" src='{base}/icons/ssh.png' alt="ssh icon">
        <span class="md:absolute  left-0 right-0  flex-1 text-center">ssh</span>
    </button>
</div>