<script lang="ts">
	import { base } from '$app/paths';
	import { goto } from '$app/navigation';
	import Button from '../components/button.svelte';

	//goto(`${base}/locked`)
</script>

<div class="flex h-full w-full items-center justify-center">
	<div class="flex-row">
		<div class="border-border_faint relative block overflow-hidden border p-3 shadow-lg">
			<div class="sm:flex sm:justify-between sm:gap-4">
				<div>
					<h3 class="text-gray-900 text-lg font-bold sm:text-xl">Enroll New Device</h3>
				</div>
			</div>

			<div class="mt-4">
				<p class="text-gray-500 text-pretty text-sm">
					Choose this option if you wish to enroll a new device with this password manager.
				</p>
			</div>

			<!-- Base -->

			<Button href="{base}/enroll">Enroll</Button>
		</div>

		<div class="border-border_faint relative mt-14 block overflow-hidden border p-3 shadow-lg">
			<div class="sm:flex sm:justify-between sm:gap-4">
				<div>
					<h3 class="text-gray-900 text-lg font-bold sm:text-xl">Register</h3>
				</div>
			</div>

			<div class="mt-4">
				<p class="text-gray-500 text-pretty text-sm">
					Choose this option to register as a first-time user with this password manager.
				</p>
			</div>

			
			<Button href="{base}/register">register</Button>
			
		</div>
	</div>
</div>
