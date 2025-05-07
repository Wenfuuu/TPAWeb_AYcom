<script lang="ts">
  import { onMount } from "svelte";
  import type { Notif } from "../model/notification";
  import type { User } from "../model/user";
  import { formatTweetDate } from "../model/thread";

    export let notifications: Notif[] = [];
    export let user: User | null;
    export let refreshData: (
        page: number,
        currTab: "all" | "mentions"
    ) => Promise<void>;

    let activeTab: "all" | "mentions" = "all";
    const page = 1;

    onMount(async() => {
        await refreshData(page, activeTab);
    });

</script>

{#if user}
    <main class="content-area">

        <!-- Header -->
        <div class="header">
            <h2>Notifications</h2>
        </div>

        <!-- Tabs -->
        <div class="tabs">
            <button
                type="button"
                class="tab-button"
                class:active={activeTab === "all"}
                on:click={() => {
                activeTab = "all";
                refreshData(page, activeTab);
                }}
            >
                All
            </button>
            <button
                type="button"
                class="tab-button"
                class:active={activeTab === "mentions"}
                on:click={() => {
                activeTab = "mentions";
                refreshData(page, activeTab);
                }}
            >
                Mentions
            </button>
        </div>

        <!-- Notifications list -->
        <div class="notifications-list">
            {#if notifications.length === 0}
                <div class="empty-notifications">
                    <p>No notifications yet</p>
                </div>
            {:else}
                {#each notifications as notification (notification.id)}
                    <div class="notification-item">
                        <div class="notification-avatar">
                            <img src={notification.profile_picture || "/placeholder.svg?height=40&width=40"} alt="Profile" />
                        </div>
                        <div class="notification-content">
                            <div class="notification-text">
                                {notification.content}
                            </div>
                            <div class="notification-time">
                                {formatTweetDate(new Date(notification.created_at))}
                            </div>
                        </div>
                    </div>
                {/each}
            {/if}
        </div>
    </main>
{/if}

<style>
    .content-area {
        flex: 1;
        border-right: 1px solid var(--border-color);
        width: 50%;
        height: 100vh;
        margin: 0 auto;
        overflow-y: auto;
        background-color: var(--background);
        color: var(--text-primary);
    }

    .header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 12px 16px;
        border-bottom: 1px solid var(--border-color);
        position: sticky;
        top: 0;
        background-color: var(--background);
        z-index: 10;
        backdrop-filter: blur(10px);
    }

    .header h2 {
        font-size: 20px;
        font-weight: 700;
        margin: 0;
    }

    .tabs {
        display: flex;
        border-bottom: 1px solid var(--border-color);
        background-color: var(--background);
        position: sticky;
        top: 53px;
        z-index: 9;
    }

    .tab-button {
        flex: 1;
        background: none;
        border: none;
        color: var(--text-tertiary);
        font-size: 15px;
        font-weight: 500;
        padding: 16px;
        cursor: pointer;
        position: relative;
        transition: color 0.2s;
    }

    .tab-button.active {
        color: var(--button-primary);
        font-weight: 700;
    }

    .tab-button.active::after {
        content: "";
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        height: 4px;
        background-color: var(--button-primary);
        border-radius: 4px 4px 0 0;
    }

    .notifications-list {
        display: flex;
        flex-direction: column;
    }

    .notification-item {
        display: flex;
        padding: 12px 16px;
        border-bottom: 1px solid var(--border-color);
        position: relative;
    }

    .notification-avatar {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        overflow: hidden;
        margin-right: 12px;
    }

    .notification-avatar img {
        width: 100%;
        height: 100%;
        object-fit: cover;
    }

    .notification-content {
        flex: 1;
        min-width: 0;
    }

    .notification-text {
        font-size: 15px;
        margin-bottom: 4px;
        word-wrap: break-word;
    }

    .notification-time {
        font-size: 14px;
        color: var(--text-tertiary);
    }

    .empty-notifications {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 200px;
        color: var(--text-tertiary);
        text-align: center;
    }

    @media (max-width: 768px) {
        .content-area {
            width: 100%;
        }
    }
</style>