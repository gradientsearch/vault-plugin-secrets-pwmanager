<script lang="ts">
	import Enroll from '../layouts/enroll.svelte';
	import Locked from '../layouts/locked.svelte';
	import Register from '../layouts/register.svelte';
	import Unlocked from '../layouts/unlocked.svelte';

	let quadrants = 3;

	let bg = ['bg-sky-500', 'bg-slate-500', 'bg-pink-400', 'bg-green-400'];
</script>

<div class="h-screen w-full text-3xl">
	<div class="flex h-full flex-1 flex-col">
		<div class="flex-1">
			<!-- <Unlocked></Unlocked> -->
			<!-- <Register></Register> -->
			<!-- <Enroll></Enroll> -->
             <Locked></Locked>
		</div>
	</div>
</div>
