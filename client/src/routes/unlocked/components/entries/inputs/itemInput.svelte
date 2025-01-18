<script lang="ts">
	import { MODE } from "../../../models/entry";


let {label, type, placeholder,  value = $bindable(), idx, last, id, mode = $bindable<MODE>()} = $props()
let disabled = $state(false)

$effect(()=>{
		mode;

		if (mode=== MODE.EDIT) {
			disabled = false
		} else {
			disabled = true
		}
	})
</script>

<label
	for="username"
	class="border-gray-200 focus-within:border-blue-600 focus-within:ring-blue-600 block overflow-hidden {idx === 0 ? 'rounded-t-md border-x border-t': ''} {last ? 'rounded-b-md border': ''}  px-3 py-2 shadow-sm focus-within:ring-1"
>
	<div class="text-gray-700 text-xs font-medium">
		{label}

		<input
			autocomplete="off"
			type="{type}"
			id="{label}"
			placeholder="{placeholder}"
			class="focus:border-transparent mt-1 w-full border-none p-0 focus:outline-none focus:ring-0 sm:text-sm"
            bind:value={value}
			disabled={mode === MODE.VIEW}
		/>
		MODE {mode === MODE.VIEW}
	</div>
</label>
