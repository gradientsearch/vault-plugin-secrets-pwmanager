<script lang="ts">
	let { className } = $props();
</script>

<div class="{className}" >
	<svg
		class="flight-icon flight-icon-vault"
		aria-hidden="true"
		data-test-icon="vault"
		fill="currentColor"
		id="icon-ember4"
		width="100%"
		height="100%"
		viewBox="0 0 24 24"
		xmlns="http://www.w3.org/2000/svg"
	>
		<path
			d="M1 1l10.96 21.334L23 1H1zm9.256 8.469H8.51V7.723h1.746v1.746zm0-2.62H8.51V5.105h1.746v1.744zm2.618 5.238h-1.746v-1.746h1.746v1.746zm0-2.618h-1.746V7.723h1.746v1.746zm0-2.62h-1.746V5.105h1.746v1.744zm2.604 2.62h-1.746V7.723h1.746v1.746zm-1.746-2.62V5.105h1.746v1.744h-1.746z"
		></path>
	</svg>
</div>

<style>
	.brand-icon-large {
		width: 62px;
		height: 62px;
	}

	.nav-header-icon {
		height: 48px;

	}
	.nav-header-icon > svg {
		@apply fill-token_side_nav_color_foreground_strong;
	}
</style>
