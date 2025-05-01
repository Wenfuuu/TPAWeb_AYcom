<script lang="ts">
  import { onMount } from "svelte";
  import { type User, type UserThread } from "../model/user";
  import {
    fetchCommunity,
    fetchThreadWithUser,
    type Community,
    type Thread,
  } from "../model/thread";

  import ComposePost from "./ComposePost.svelte";

  export let threads: (Thread & { userData: UserThread })[] = [];
  let communities: Community[] = [];

  export let refreshData: (page: number) => Promise<void>;

  export let user: User | null;

  let page = 1;
  let isLoading = false;
  let hasMore = true;

  let contentAreaElement: HTMLElement | null = null;

  let sentinel: HTMLDivElement | null = null;
  let observer: IntersectionObserver | null = null;

  let currentActiveTab: "for-you" | "following" = "for-you";

  function switchTab(tab: "for-you" | "following") {
    if (currentActiveTab !== tab) {
      currentActiveTab = tab;
      page = 1;
      threads = [];
      hasMore = true;

      refreshData(page);
    }
  }

  function observeScroll() {
    if (observer) {
      observer.disconnect();
    }

    if (!contentAreaElement) {
      // Retry after the current event loop tick, allowing bind:this to complete
      requestAnimationFrame(observeScroll);
      return;
    }

    observer = new IntersectionObserver(
      async (entries) => {
        console.log(hasMore);
        if (entries[0].isIntersecting && !isLoading && hasMore) {
          isLoading = true;
          console.log("Loading more content, page:", page + 1);

          const newThreads = await fetchThreadWithUser(page + 1);
          if(!newThreads){
            hasMore = false;
            isLoading = false;
            return;
          }
          console.log("Fetched new threads:", newThreads.length);
          threads = [...threads, ...newThreads];
          page += 1;

          isLoading = false;
        }
      },
      {
        root: contentAreaElement,
        rootMargin: "1px",
        threshold: 1,
      }
    );

    if (sentinel) {
      observer.observe(sentinel);
      console.log(
        "Observer attached to sentinel element with content area root."
      );
    } else {
      console.warn(
        "Sentinel element not found - scroll observation not active"
      );
    }
  }

  onMount(async () => {
    observeScroll();
    await refreshData(page);
    communities = await fetchCommunity();
  });
</script>

{#if user}
  <main class="content-area" bind:this={contentAreaElement}>
    <div class="header">
      <div class="tabs">
        <button
          type="button"
          class="tab"
          class:active={currentActiveTab === "for-you"}
          on:click={() => switchTab("for-you")}
        >
          For you
        </button>
        <button
          type="button"
          class="tab"
          class:active={currentActiveTab === "following"}
          on:click={() => switchTab("following")}
        >
          Following
        </button>
      </div>
    </div>

    <ComposePost {user} {refreshData} />

    <!-- Feed Display Section -->
    <div class="posts">
      {#each threads as thread (thread.id)}
        <div class="post">
          <div class="post-avatar">
            <!-- Use optional chaining and provide a fallback -->
            <img
              src={thread.userData?.profile || "/placeholder.svg"}
              alt={thread.userData?.name || "User"}
            />
          </div>
          <div class="post-content">
            <div class="post-header">
              <div class="post-user-info">
                <!-- Use optional chaining for user data -->
                <span class="post-user-name"
                  >{thread.userData?.name || "Unknown User"}</span
                >
                {#if thread.userData?.isPremium}
                  <span class="verified-badge">
                    <svg
                      viewBox="0 0 24 24"
                      width="16"
                      height="16"
                      fill="rgb(29, 155, 240)"
                    >
                      <path
                        d="M22.25 12c0-1.43-.88-2.67-2.19-3.34.46-1.39.2-2.9-.81-3.91s-2.52-1.27-3.91-.81c-.66-1.31-1.91-2.19-3.34-2.19s-2.67.88-3.33 2.19c-1.4-.46-2.91-.2-3.92.81s-1.26 2.52-.8 3.91c-1.31.67-2.2 1.91-2.2 3.34s.89 2.67 2.2 3.34c-.46 1.39-.21 2.9.8 3.91s2.52 1.26 3.91.81c.67 1.31 1.91 2.19 3.34 2.19s2.68-.88 3.34-2.19c1.39.45 2.9.2 3.91-.81s1.27-2.52.81-3.91c1.31-.67 2.19-1.91 2.19-3.34zm-11.71 4.2L6.8 12.46l1.41-1.42 2.26 2.26 4.8-5.23 1.47 1.36-6.2 6.77z"
                      ></path>
                    </svg>
                  </span>
                {/if}
                <span class="post-username"
                  >@{thread.userData?.username || "unknown"}</span
                >
                <!-- Format date nicely later if needed -->
                <span class="post-date">· {thread.created_at}</span>
              </div>
              <button type="button" class="post-more" aria-label="Post More">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <circle cx="12" cy="12" r="1"></circle>
                  <circle cx="19" cy="12" r="1"></circle>
                  <circle cx="5" cy="12" r="1"></circle>
                </svg>
              </button>
            </div>

            {#if thread.community_id != "00000000-0000-0000-0000-000000000000"}
              <div class="post-community">
                Community: {communities.find(
                  (c) => c.id === thread.community_id
                )?.name || "Community"}
              </div>
            {/if}

            {#if thread.content}
              <div class="post-text">{thread.content}</div>
            {/if}

            {#if thread.images && thread.images.length > 0}
              <div
                class="post-images"
                class:grid-2={thread.images.length === 2}
                class:grid-3={thread.images.length === 3}
                class:grid-4={thread.images.length === 4}
              >
                {#each thread.images as image, i (i)}
                  <div class="image-container">
                    <img
                      src={image || "/placeholder.svg"}
                      alt="Post image {i + 1}"
                    />
                  </div>
                {/each}
              </div>
            {/if}

            <div class="post-actions">
              <button type="button" class="post-action">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path
                    d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"
                  ></path>
                </svg>
                <span>{thread.insights?.replies || 0}</span>
              </button>
              <button type="button" class="post-action">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="m17 2 4 4-4 4"></path>
                  <path d="M3 11v-1a4 4 0 0 1 4-4h14"></path>
                  <path d="m7 22-4-4 4-4"></path>
                  <path d="M21 13v1a4 4 0 0 1-4 4H3"></path>
                </svg>
                <span>{thread.insights?.reposts || 0}</span>
              </button>
              <button type="button" class="post-action">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path
                    d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"
                  ></path>
                </svg>
                <!-- Display likes directly, handle large numbers formatting separately if needed -->
                <span>{thread.insights?.likes || 0}</span>
              </button>
              <button type="button" class="post-action" aria-label="Share Post">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                  <polyline points="17 8 12 3 7 8"></polyline>
                  <line x1="12" x2="12" y1="3" y2="15"></line>
                </svg>
              </button>
            </div>
          </div>
        </div>
      {/each}
    </div>

    <div bind:this={sentinel} class="loading-sentinel">
      {#if isLoading}
        <!-- Skeleton Loading UI -->
        <div class="skeleton-loading">
          <!-- Repeat skeleton post structure multiple times for better UX -->
          {#each Array(3) as _, i (i)}
            <div class="skeleton-post">
              <div class="skeleton-avatar pulse"></div>
              <div class="skeleton-content">
                <div class="skeleton-header">
                  <div class="skeleton-name pulse"></div>
                  <div class="skeleton-username pulse"></div>
                </div>
                <div class="skeleton-text pulse"></div>
                <div class="skeleton-text-short pulse"></div>
                <!-- Random image placeholder in some posts -->
                {#if i % 2 === 0}
                  <div class="skeleton-image pulse"></div>
                {/if}
                <div class="skeleton-actions">
                  <div class="skeleton-action pulse"></div>
                  <div class="skeleton-action pulse"></div>
                  <div class="skeleton-action pulse"></div>
                  <div class="skeleton-action pulse"></div>
                </div>
              </div>
            </div>
          {/each}
        </div>
      {:else if !hasMore && threads.length > 0}
        <p>No more thread available!</p>
      {/if}
    </div>
  </main>
{/if}

<style>
  /* Styles for HomeContent - Keep only styles relevant to this component's structure */
  .content-area {
    flex: 1;
    border-right: 1px solid var(--border-color);
    width: 50%;
    height: 100vh;
    margin: 0 auto; /* Center if using max-width */
    overflow-y: auto;
  }

  .header {
    position: sticky;
    top: 0;
    background-color: var(--card-background);
    backdrop-filter: blur(12px);
    z-index: 10;
    border-bottom: 1px solid var(--border-color);
  }

  .tabs {
    display: flex;
    height: 53px;
  }

  .tab {
    flex: 1;
    background: none;
    border: none;
    color: var(--text-secondary); /* Default color */
    font-size: 15px;
    font-weight: 500;
    padding: 16px;
    cursor: pointer;
    position: relative;
    transition:
      background-color 0.2s,
      color 0.2s;
  }

  .tab:hover {
    background-color: var(--hover-background);
    color: var(--text-primary);
  }

  .tab.active {
    color: var(--text-primary);
    font-weight: 700;
  }

  .tab.active::after {
    content: "";
    position: absolute;
    bottom: 0;
    left: 50%;
    transform: translateX(-50%);
    width: 56px;
    height: 4px;
    background-color: var(--button-primary);
    border-radius: 9999px;
  }

  /* Styles for the posts list */
  .posts {
    display: flex;
    flex-direction: column;
  }

  .post {
    display: flex;
    padding: 12px 16px;
    border-bottom: 1px solid var(--border-color);
    transition: background-color 0.2s;
    cursor: pointer; /* Indicate posts are clickable */
  }

  .post:hover {
    background-color: var(--hover-background);
  }

  .post-avatar {
    margin-right: 12px;
    flex-shrink: 0; /* Prevent avatar shrinking */
  }

  .post-avatar img {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover; /* Ensure image covers the area */
  }

  .post-content {
    flex: 1;
    min-width: 0; /* Prevent content overflow issues */
  }

  .post-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start; /* Align items to top */
    margin-bottom: 4px;
  }

  .post-user-info {
    display: flex;
    align-items: center;
    flex-wrap: wrap; /* Allow wrapping on small screens */
    gap: 4px; /* Space between elements */
    min-width: 0;
    line-height: 1.3; /* Adjust line height */
  }

  .post-user-name {
    font-weight: 700;
    color: var(--text-primary);
    white-space: nowrap;
    overflow: hidden; /* Prevent long names breaking layout */
    text-overflow: ellipsis;
  }

  .verified-badge {
    display: flex; /* Use flex for alignment */
    align-items: center;
    flex-shrink: 0; /* Prevent shrinking */
  }

  .post-username,
  .post-date {
    color: var(--text-tertiary);
    white-space: nowrap;
    font-size: 14px; /* Slightly smaller */
  }

  .post-username {
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .post-more {
    background: none;
    border: none;
    color: var(--text-tertiary);
    cursor: pointer;
    padding: 4px; /* Easier clicking */
    margin-left: 8px; /* Space from user info */
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition:
      background-color 0.2s,
      color 0.2s;
    width: 30px; /* Explicit size */
    height: 30px;
    flex-shrink: 0;
  }

  .post-more:hover {
    background-color: rgba(var(--button-primary-rgb, 29, 155, 240), 0.1);
    color: var(--button-primary);
  }

  .post-community {
    color: var(--button-primary);
    font-size: 13px; /* Smaller font size */
    margin-bottom: 8px;
    font-weight: 500;
  }

  .post-text {
    color: var(--text-primary);
    margin-bottom: 12px;
    word-wrap: break-word; /* Ensure text wraps */
    line-height: 1.5; /* Improve readability */
    font-size: 15px;
  }

  .post-images {
    border-radius: 16px;
    overflow: hidden;
    margin-top: 8px; /* Space between text and images */
    margin-bottom: 12px;
    display: grid;
    gap: 2px;
    border: 1px solid var(--border-color); /* Subtle border */
  }

  /* Grid layouts for different image counts */
  .post-images.grid-2 {
    grid-template-columns: repeat(2, 1fr);
    aspect-ratio: 16 / 9;
  }
  .post-images.grid-3 {
    grid-template-columns: repeat(2, 1fr);
    grid-template-rows: repeat(2, 1fr);
    aspect-ratio: 1 / 1;
  }
  .post-images.grid-3 .image-container:first-child {
    grid-column: 1 / 2;
    grid-row: 1 / 3;
  }
  .post-images.grid-3 .image-container:nth-child(2) {
    grid-column: 2 / 3;
    grid-row: 1 / 2;
  }
  .post-images.grid-3 .image-container:nth-child(3) {
    grid-column: 2 / 3;
    grid-row: 2 / 3;
  }
  .post-images.grid-4 {
    grid-template-columns: repeat(2, 1fr);
    grid-template-rows: repeat(2, 1fr);
    aspect-ratio: 1 / 1;
  }

  .image-container {
    position: relative; /* For potential future overlays */
    overflow: hidden;
    background-color: var(--hover-background); /* Placeholder background */
  }

  .image-container img {
    display: block; /* Remove bottom space */
    width: 100%;
    height: 100%;
    object-fit: cover; /* Cover the area */
    transition: opacity 0.3s ease; /* Smooth loading */
  }

  .post-actions {
    display: flex;
    justify-content: space-between;
    max-width: 425px; /* Limit width of actions row */
    margin-top: 12px;
  }

  .post-action {
    background: none;
    border: none;
    color: var(--text-tertiary);
    display: flex;
    align-items: center;
    gap: 6px; /* Increased gap */
    cursor: pointer;
    transition: color 0.2s;
    font-size: 13px;
    padding: 6px; /* Make touch target larger */
    border-radius: 4px; /* Slight rounding on hover */
  }

  .post-action:hover {
    color: var(--button-primary);
    background-color: rgba(var(--button-primary-rgb, 29, 155, 240), 0.1);
  }

  .post-action svg {
    width: 18px;
    height: 18px;
  }

  .loading-sentinel {
    min-height: 100px;
    padding: 20px;
    text-align: center;
    color: var(--text-secondary);
  }

  /* Skeleton loading styles */
  .skeleton-loading {
    width: 100%;
    padding: 0 16px;
  }

  .skeleton-post {
    display: flex;
    padding: 12px 0;
    border-bottom: 1px solid var(--border-color);
    animation: fadeIn 0.5s ease-in-out;
  }

  .skeleton-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background-color: var(--hover-background);
    margin-right: 12px;
    flex-shrink: 0;
  }

  .skeleton-content {
    flex: 1;
  }

  .skeleton-header {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
    gap: 8px;
  }

  .skeleton-name {
    height: 16px;
    width: 120px;
    background-color: var(--hover-background);
    border-radius: 4px;
  }

  .skeleton-username {
    height: 14px;
    width: 80px;
    background-color: var(--hover-background);
    border-radius: 4px;
  }

  .skeleton-text {
    height: 14px;
    width: 100%;
    background-color: var(--hover-background);
    border-radius: 4px;
    margin-bottom: 8px;
  }

  .skeleton-text-short {
    height: 14px;
    width: 70%;
    background-color: var(--hover-background);
    border-radius: 4px;
    margin-bottom: 16px;
  }

  .skeleton-image {
    height: 180px;
    width: 100%;
    background-color: var(--hover-background);
    border-radius: 16px;
    margin-bottom: 16px;
  }

  .skeleton-actions {
    display: flex;
    justify-content: space-between;
    max-width: 425px;
  }

  .skeleton-action {
    height: 18px;
    width: 36px;
    background-color: var(--hover-background);
    border-radius: 4px;
  }

  @media (max-width: 768px) {
    .content-area {
      border-right: none; /* Remove border on mobile */
    }
    .post-actions {
      max-width: 100%; /* Allow actions to take full width */
    }
  }

  @media (max-width: 500px) {
    .post-user-info {
      font-size: 14px; /* Slightly smaller text on very small screens */
    }
    .post-text {
      font-size: 14px;
    }
    .post-action span {
      display: none;
    }
    .post-action {
      gap: 4px;
    }
  }
</style>
