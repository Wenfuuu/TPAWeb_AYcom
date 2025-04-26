<script lang="ts">
    const { text, type, disabled = false, onClick } = $props();
</script>

<button
    onclick={onClick}
    {disabled}
    {type}
>
    {text}
</button>

<style>

    button {
        background-color: rgb(88, 88, 88);
        color: white;
        border: none;
        flex: 1;
        padding: 12px;
        border-radius: 9999px;
        font-size: 15px;
        font-weight: 600;
        cursor: pointer;
        transition: background-color 0.2s;
    }

    button:hover:not(:disabled) {
        background-color: rgba(255, 255, 255, 0.1);
    }

    button:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

</style>