<script lang="ts">
  import { onMount } from "svelte";
  import {
    bookmarkThread,
    fetchThreadRepliesWithUser,
    likeThread,
    repostThread,
    unbookmarkThread,
    unlikeThread,
    unrepostThread,
    type Thread,
  } from "../model/thread";
  import {
    fetchUserEmailById,
    type User,
    type UserThread,
  } from "../model/user";
  import { navigate } from "svelte-routing";
  import ComposeReply from "./ComposeReply.svelte";
  import { sendThreadNotif, setupWS } from "../model/websocket";

  export let thread: Thread & { userData: UserThread };
  export let replies: (Thread & { userData: UserThread })[] = [];

  export let refreshData: (page: number) => Promise<void>;
  export let loadThreadData: (threadId: string, page: number) => Promise<void>;
  export let handleBackNavigation: () => void;
  //   export let changePrevId: () => void;

  export let user: User | null;
  export let socket: WebSocket;
  let page = 1;
  let isLoading = false;
  let hasMore = true;
  let contentAreaElement: HTMLElement | null = null;

  let sentinel: HTMLDivElement | null = null;
  let observer: IntersectionObserver | null = null;
  // Media overlay state
  let showMediaOverlay = false;
  let currentMediaIndex = 0;
  let overlayMediaItems: string[] = [];

  // Reply composition state
  let showReplyTextarea = false;
  let replyText = "";

  function observeScroll() {
    if (observer) {
      observer.disconnect();
    }

    if (!contentAreaElement) {
      requestAnimationFrame(observeScroll);
      return;
    }

    observer = new IntersectionObserver(
      async (entries) => {
        // console.log(hasMore);
        if (entries[0].isIntersecting && !isLoading && hasMore) {
          isLoading = true;
          // console.log("Loading more content, page:", page + 1);

          await new Promise((resolve) => setTimeout(resolve, 500));

          let newReplies;
          if (user) {
            newReplies = await fetchThreadRepliesWithUser(
              page + 1,
              thread.id,
              user.id
            );
          }

          if (!newReplies) {
            hasMore = false;
            isLoading = false;
            return;
          }
          // console.log("Fetched new replies:", newReplies.length);
          replies = [...replies, ...newReplies];
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
      // console.log(
      //   "Observer attached to sentinel element with content area root."
      // );
    } else {
      // console.warn(
      //   "Sentinel element not found - scroll observation not active"
      // );
    }
  }
  async function handleReplyCreated() {
    hasMore = true;
    loadThreadData(thread.id, page);
  }

  // Media overlay functions
  function openMediaOverlay(mediaItems: string[], startIndex: number) {
    overlayMediaItems = mediaItems;
    currentMediaIndex = startIndex;
    showMediaOverlay = true;
    document.body.style.overflow = "hidden"; // Prevent background scrolling
  }

  function closeMediaOverlay() {
    showMediaOverlay = false;
    document.body.style.overflow = ""; // Restore scrolling
  }

  function navigateToNextMedia() {
    if (currentMediaIndex < overlayMediaItems.length - 1) {
      currentMediaIndex++;
    }
  }

  function navigateToPrevMedia() {
    if (currentMediaIndex > 0) {
      currentMediaIndex--;
    }
  }

  function handleOverlayKeydown(event: KeyboardEvent) {
    if (event.key === "Escape") {
      closeMediaOverlay();
    } else if (event.key === "ArrowRight") {
      navigateToNextMedia();
    } else if (event.key === "ArrowLeft") {
      navigateToPrevMedia();
    }
  }
  function toggleLike(targetThread: Thread) {
    targetThread.isLiked = !targetThread.isLiked;
  }

  function toggleRepost(targetThread: Thread) {
    targetThread.isReposted = !targetThread.isReposted;
  }

  async function handleLikeThread(): Promise<void> {
    if (!user) return;

    if (!thread.isLiked) {
      await likeThread(user.id, thread.id);

      const email = await fetchUserEmailById(thread.user_id, user.id);
      sendThreadNotif(
        socket,
        thread.user_id,
        user.id,
        thread.id,
        `You received a like from ${user.name}`,
        email
      );
    } else {
      await unlikeThread(user.id, thread.id);
    }

    toggleLike(thread);
    thread = { ...thread };
  }

  async function handleLikeReplies(replyThread: Thread): Promise<void> {
    if (!user) return;

    if (!replyThread.isLiked) {
      await likeThread(user.id, replyThread.id);

      const email = await fetchUserEmailById(replyThread.user_id, user.id);
      sendThreadNotif(
        socket,
        replyThread.user_id,
        user.id,
        thread.id,
        `You received a like from ${user.name}`,
        email
      );
    } else {
      await unlikeThread(user.id, replyThread.id);
    }

    toggleLike(replyThread);
    replies = [...replies];
  }

  function toggleBookmark(targetThread: Thread) {
    targetThread.isBookmarked = !targetThread.isBookmarked;
  }

  async function handleBookmarkThread(): Promise<void> {
    if (!user) return;

    if (!thread.isBookmarked) {
      await bookmarkThread(user.id, thread.id);

      const email = await fetchUserEmailById(thread.user_id, user.id);
      sendThreadNotif(
        socket,
        thread.user_id,
        user.id,
        thread.id,
        `You received a bookmark from ${user.name}`,
        email
      );
    } else {
      await unbookmarkThread(user.id, thread.id);
    }

    toggleBookmark(thread);
    thread = { ...thread };
  }
  async function handleBookmarkReplies(replyThread: Thread): Promise<void> {
    if (!user) return;

    if (!replyThread.isBookmarked) {
      await bookmarkThread(user.id, replyThread.id);

      const email = await fetchUserEmailById(replyThread.user_id, user.id);
      sendThreadNotif(
        socket,
        thread.user_id,
        user.id,
        thread.id,
        `You received a bookmark from ${user.name}`,
        email
      );
    } else {
      await unbookmarkThread(user.id, replyThread.id);
    }

    toggleBookmark(replyThread);
    replies = [...replies];
  }
  async function handleRepostThread(): Promise<void> {
    if (!user) return;

    if (!thread.isReposted) {
      await repostThread(user.id, thread.id);

      const email = await fetchUserEmailById(thread.user_id, user.id);
      sendThreadNotif(
        socket,
        thread.user_id,
        user.id,
        thread.id,
        `You received a repost from ${user.name}`,
        email
      );
    } else {
      await unrepostThread(user.id, thread.id);
    }

    toggleRepost(thread);
    thread = { ...thread };
  }

  async function handleRepostReplies(replyThread: Thread): Promise<void> {
    if (!user) return;

    if (!replyThread.isReposted) {
      await repostThread(user.id, replyThread.id);

      const email = await fetchUserEmailById(replyThread.user_id, user.id);
      sendThreadNotif(
        socket,
        replyThread.user_id,
        user.id,
        replyThread.id,
        `You received a repost from ${user.name}`,
        email
      );
    } else {
      await unrepostThread(user.id, replyThread.id);
    }

    toggleRepost(replyThread);
    replies = [...replies];
  }

  // Reply submission handler
  async function handleReplySubmit(): Promise<void> {
    if (!replyText.trim() || !user) return;

    try {
      // TODO: Implement reply submission API call

      // Reset form
      replyText = "";
      showReplyTextarea = false;

      // Refresh replies
      await refreshData(1);
    } catch {
      // TODO: Handle error properly
    }
  }

  onMount(async () => {
    observeScroll();
    if (thread) loadThreadData(thread.id, 1);
    if (!user) return;
    socket = setupWS(socket, user.id);
  });
</script>

{#if user && thread}
  <div class="thread-detail-container">
    <!-- Left side: Main content -->
    <main class="content-area" bind:this={contentAreaElement}>
      <div class="header">
        <button
          type="button"
          class="back-button"
          aria-label="Back button"
          on:click={() => {
            handleBackNavigation();
          }}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="20"
            height="20"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="m15 18-6-6 6-6"></path>
          </svg>
        </button>
        <h1>Post</h1>
      </div>

      <div class="thread-post">
        <div class="post-avatar">
          <div
            role="button"
            tabindex="0"
            class="avatar-button"
            on:click|stopPropagation={() =>
              navigate(`/profile/${thread.userData?.user_id}`)}
            on:keydown|stopPropagation={(e) => {
              if (e.key === "Enter" || e.key === " ") {
                navigate(`/profile/${thread.userData?.user_id}`);
              }
            }}
          >
            <img
              src={thread.userData?.profile || "/placeholder.svg"}
              alt={thread.userData?.name || "User"}
            />
          </div>
        </div>
        <div class="post-content">
          <div class="post-header">
            <div class="post-user-info">
              <div class="user-name-container">
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
              </div>
              <span class="post-username"
                >@{thread.userData?.username || "unknown"}</span
              >
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
                  {#if image.endsWith(".png") || image.endsWith(".jpg") || image.endsWith(".jpeg") || image.endsWith(".gif")}
                    <div
                      role="button"
                      tabindex="0"
                      class="media-button"
                      on:click|stopPropagation={() =>
                        openMediaOverlay(thread.images, i)}
                      on:keydown={(e) => {
                        if (e.key === "Enter" || e.key === " ") {
                          e.preventDefault();
                          openMediaOverlay(thread.images, i);
                        }
                      }}
                    >
                      <img
                        src={image || "/placeholder.svg"}
                        alt="Post image {i + 1}"
                      />
                    </div>
                  {:else}
                    <div
                      role="button"
                      tabindex="0"
                      class="media-button"
                      on:click|stopPropagation={() =>
                        openMediaOverlay(thread.images, i)}
                      on:keydown={(e) => {
                        if (e.key === "Enter" || e.key === " ") {
                          e.preventDefault();
                          openMediaOverlay(thread.images, i);
                        }
                      }}
                    >
                      <video src={image} controls class="post-video">
                        <track
                          kind="captions"
                          src="path_to_captions.vtt"
                          srclang="en"
                          label="English"
                        />
                      </video>
                    </div>
                  {/if}
                </div>
              {/each}
            </div>
          {/if}

          <div class="post-time">{thread.created_at}</div>

          <div class="post-actions">
            <button
              type="button"
              class="post-action reply"
              class:active={false}
            >
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
            <button
              type="button"
              class="post-action repost"
              class:active={thread.isReposted}
              on:click|stopPropagation={handleRepostThread}
            >
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
            <button
              type="button"
              class="post-action like"
              class:active={thread.isLiked}
              on:click|stopPropagation={handleLikeThread}
            >
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
              <span>{thread.insights?.likes || 0}</span>
            </button>
            <button
              type="button"
              class="post-action bookmark"
              class:active={thread.isBookmarked}
              on:click|stopPropagation={handleBookmarkThread}
              aria-label="Bookmark Post"
            >
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
                <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v16z"
                ></path>
              </svg>
            </button>
          </div>
        </div>
      </div>

      <ComposeReply {thread} {user} refreshData={handleReplyCreated} />

      <div class="posts">
        {#each replies as replyThread (replyThread.id)}
          <div
            class="post"
            role="button"
            tabindex="0"
            on:click={() => {
              navigate(`/thread-detail/${replyThread.id}`);
              refreshData(1);
            }}
            on:keydown={(e) => {
              if (e.key === "Enter" || e.key === " ") {
                navigate(`/thread-detail/${replyThread.id}`);
              }
            }}
          >
            <div class="post-avatar">
              <div
                role="button"
                tabindex="0"
                class="avatar-button"
                on:click|stopPropagation={() =>
                  navigate(`/profile/${replyThread.userData?.user_id}`)}
                on:keydown|stopPropagation={(e) => {
                  if (e.key === "Enter" || e.key === " ") {
                    navigate(`/profile/${replyThread.userData?.user_id}`);
                  }
                }}
              >
                <img
                  src={replyThread.userData?.profile || "/placeholder.svg"}
                  alt={replyThread.userData?.name || "User"}
                />
              </div>
            </div>
            <div class="post-content">
              <div class="reply-to">
                Replied to <span class="reply-username"
                  >@{thread.userData?.username || "unknown"}</span
                >
              </div>
              <div class="post-header">
                <div class="post-reply-info">
                  <span class="post-user-name"
                    >{replyThread.userData?.name || "Unknown User"}</span
                  >
                  {#if replyThread.userData?.isPremium}
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
                    >@{replyThread.userData?.username || "unknown"}</span
                  >
                  <span class="post-date">· {replyThread.created_at}</span>
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

              {#if replyThread.content}
                <div class="post-text">{replyThread.content}</div>
              {/if}

              {#if replyThread.images && replyThread.images.length > 0}
                <div
                  class="post-images"
                  class:grid-2={replyThread.images.length === 2}
                  class:grid-3={replyThread.images.length === 3}
                  class:grid-4={replyThread.images.length === 4}
                >
                  {#each replyThread.images as image, i (i)}
                    <div class="image-container">
                      <div
                        role="button"
                        tabindex="0"
                        class="media-button"
                        on:click|stopPropagation={() =>
                          openMediaOverlay(replyThread.images, i)}
                        on:keydown={(e) => {
                          if (e.key === "Enter" || e.key === " ") {
                            e.preventDefault();
                            openMediaOverlay(replyThread.images, i);
                          }
                        }}
                      >
                        <img
                          src={image || "/placeholder.svg"}
                          alt="Post image {i + 1}"
                        />
                      </div>
                    </div>
                  {/each}
                </div>
              {/if}

              <div class="post-actions">
                <button
                  type="button"
                  class="post-action reply"
                  class:active={false}
                >
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
                  <span>{replyThread.insights?.replies || 0}</span>
                </button>
                <button
                  type="button"
                  class="post-action repost"
                  class:active={replyThread.isReposted}
                  on:click|stopPropagation={() =>
                    handleRepostReplies(replyThread)}
                >
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
                  <span>{replyThread.insights?.reposts || 0}</span>
                </button>
                <button
                  type="button"
                  class="post-action like"
                  class:active={replyThread.isLiked}
                  on:click|stopPropagation={() =>
                    handleLikeReplies(replyThread)}
                >
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
                  <span>{replyThread.insights?.likes || 0}</span>
                </button>
                <button
                  type="button"
                  class="post-action bookmark"
                  class:active={replyThread.isBookmarked}
                  on:click|stopPropagation={() =>
                    handleBookmarkReplies(replyThread)}
                  aria-label="Bookmark Post"
                >
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
                      d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v16z"
                    ></path>
                  </svg>
                </button>
              </div>
            </div>
          </div>
        {/each}
      </div>

      <div bind:this={sentinel} class="loading-sentinel">
        {#if isLoading}
          <div class="skeleton-loading">
            <!-- eslint-disable-next-line @typescript-eslint/no-unused-vars -->
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
        {:else if !hasMore && replies.length > 0}
          <p>No more replies available!</p>
        {:else if replies.length === 0 && !isLoading}
          <div class="no-replies">
            <p>No replies yet. Be the first to reply!</p>
          </div>
        {/if}
      </div>
    </main>
  </div>

  <!-- Media Overlay -->
  {#if showMediaOverlay}
    <div
      class="media-overlay"
      role="dialog"
      aria-modal="true"
      aria-label="Media viewer"
      on:keydown={handleOverlayKeydown}
      tabindex="-1"
    >
      <div
        class="media-overlay-backdrop"
        role="button"
        tabindex="0"
        on:click={closeMediaOverlay}
        on:keydown={(e) => {
          if (e.key === "Enter" || e.key === " ") {
            e.preventDefault();
            closeMediaOverlay();
          }
        }}
      ></div>

      <div class="media-overlay-content">
        <!-- Close button -->
        <button
          type="button"
          class="overlay-close-button"
          aria-label="Close media overlay"
          on:click={closeMediaOverlay}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="m18 6-12 12"></path>
            <path d="m6 6 12 12"></path>
          </svg>
        </button>

        <!-- Left side: Media viewer -->
        <div class="media-overlay-left">
          <div class="media-overlay-viewer">
            <!-- Navigation arrows -->
            {#if overlayMediaItems.length > 1}
              <button
                type="button"
                class="media-nav-button prev"
                class:disabled={currentMediaIndex === 0}
                aria-label="Previous media"
                on:click={navigateToPrevMedia}
                disabled={currentMediaIndex === 0}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="m15 18-6-6 6-6"></path>
                </svg>
              </button>

              <button
                type="button"
                class="media-nav-button next"
                class:disabled={currentMediaIndex ===
                  overlayMediaItems.length - 1}
                aria-label="Next media"
                on:click={navigateToNextMedia}
                disabled={currentMediaIndex === overlayMediaItems.length - 1}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="m9 18 6-6-6-6"></path>
                </svg>
              </button>
            {/if}

            <!-- Media content -->
            <div class="media-overlay-media">
              {#if overlayMediaItems[currentMediaIndex]}
                {#if overlayMediaItems[currentMediaIndex].endsWith(".png") || overlayMediaItems[currentMediaIndex].endsWith(".jpg") || overlayMediaItems[currentMediaIndex].endsWith(".jpeg") || overlayMediaItems[currentMediaIndex].endsWith(".gif")}
                  <img
                    src={overlayMediaItems[currentMediaIndex]}
                    alt="Media {currentMediaIndex +
                      1} of {overlayMediaItems.length}"
                    class="overlay-media-item"
                  />
                {:else}
                  <video
                    src={overlayMediaItems[currentMediaIndex]}
                    controls
                    class="overlay-media-item overlay-video"
                    autoplay
                  >
                    <track
                      kind="captions"
                      src="path_to_captions.vtt"
                      srclang="en"
                      label="English"
                    />
                  </video>
                {/if}
              {/if}
            </div>

            <!-- Media counter -->
            {#if overlayMediaItems.length > 1}
              <div class="media-counter">
                {currentMediaIndex + 1} / {overlayMediaItems.length}
              </div>
            {/if}
          </div>
        </div>

        <!-- Right side: Thread interactions -->
        <div class="media-overlay-right">
          <div class="overlay-sidebar">
            <!-- Original thread -->
            <div class="overlay-original-thread">
              <div class="overlay-thread-header">
                <div class="post-avatar">
                  <img
                    src={thread.userData?.profile || "/placeholder.svg"}
                    alt={thread.userData?.name || "User"}
                  />
                </div>
                <div class="overlay-user-info">
                  <div class="user-name-container">
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
                  </div>
                  <span class="post-username"
                    >@{thread.userData?.username || "unknown"}</span
                  >
                  <span class="post-date">· {thread.created_at}</span>
                </div>
              </div>

              {#if thread.content}
                <div class="overlay-post-text">{thread.content}</div>
              {/if}

              <div class="overlay-post-actions">
                <button
                  type="button"
                  class="post-action reply"
                  class:active={false}
                  on:click={() => (showReplyTextarea = !showReplyTextarea)}
                >
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
                <button
                  type="button"
                  class="post-action repost"
                  on:click|stopPropagation={handleRepostThread}
                >
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
                <button
                  type="button"
                  class="post-action like"
                  class:active={thread.isLiked}
                  on:click|stopPropagation={handleLikeThread}
                >
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
                  <span>{thread.insights?.likes || 0}</span>
                </button>
                <button
                  type="button"
                  class="post-action bookmark"
                  class:active={thread.isBookmarked}
                  on:click|stopPropagation={handleBookmarkThread}
                  aria-label="Bookmark Post"
                >
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
                      d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v16z"
                    ></path>
                  </svg>
                </button>
              </div>
            </div>

            <!-- Reply composition area -->
            {#if showReplyTextarea}
              <div class="overlay-reply-compose">
                <div class="reply-compose-header">
                  <div class="post-avatar">
                    <img src="/placeholder.svg" alt="Your avatar" />
                  </div>
                  <textarea
                    bind:value={replyText}
                    placeholder="Post your reply"
                    class="reply-textarea"
                    on:input={(e) => {
                      const target = e.target as HTMLTextAreaElement;
                      if (target) {
                        target.style.height = "auto";
                        target.style.height = target.scrollHeight + "px";
                      }
                    }}
                  ></textarea>
                </div>
                <div class="reply-compose-actions">
                  <button
                    type="button"
                    class="reply-submit-btn"
                    disabled={!replyText.trim()}
                    on:click={handleReplySubmit}
                  >
                    Reply
                  </button>
                </div>
              </div>
            {/if}

            <!-- ComposeReply in media overlay -->
            <div class="overlay-compose-reply">
              <h3>Reply to this thread</h3>
              <ComposeReply {thread} {user} refreshData={handleReplyCreated} />
            </div>

            <!-- Replies list -->
            <div class="overlay-replies">
              {#each replies as replyThread (replyThread.id)}
                <div class="overlay-reply-item">
                  <div class="reply-header">
                    <div class="post-avatar">
                      <img
                        src={replyThread.userData?.profile ||
                          "/placeholder.svg"}
                        alt={replyThread.userData?.name || "User"}
                      />
                    </div>
                    <div class="reply-user-info">
                      <div class="user-name-container">
                        <span class="post-user-name"
                          >{replyThread.userData?.name || "Unknown User"}</span
                        >
                        {#if replyThread.userData?.isPremium}
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
                      </div>
                      <span class="post-username"
                        >@{replyThread.userData?.username || "unknown"}</span
                      >
                      <span class="post-date">· {replyThread.created_at}</span>
                    </div>
                  </div>

                  {#if replyThread.content}
                    <div class="reply-text">{replyThread.content}</div>
                  {/if}
                  <!-- Reply images -->
                  {#if replyThread.images && replyThread.images.length > 0}
                    <div class="reply-images">
                      {#each replyThread.images as image, i (i)}
                        <div
                          role="button"
                          tabindex="0"
                          class="reply-image-button"
                          on:click={() =>
                            openMediaOverlay(replyThread.images, i)}
                          on:keydown={(e) => {
                            if (e.key === "Enter" || e.key === " ") {
                              e.preventDefault();
                              openMediaOverlay(replyThread.images, i);
                            }
                          }}
                        >
                          <img
                            src={image}
                            alt="Reply image {i + 1}"
                            class="reply-image"
                          />
                        </div>
                      {/each}
                    </div>
                  {/if}

                  <div class="reply-actions">
                    <button
                      type="button"
                      class="post-action reply"
                      class:active={false}
                    >
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
                      <span>{replyThread.insights?.replies || 0}</span>
                    </button>
                    <button
                      type="button"
                      class="post-action repost"
                      class:active={replyThread.isReposted}
                      on:click|stopPropagation={() =>
                        handleRepostReplies(replyThread)}
                    >
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
                      <span>{replyThread.insights?.reposts || 0}</span>
                    </button>
                    <button
                      type="button"
                      class="post-action like"
                      class:active={replyThread.isLiked}
                      on:click|stopPropagation={() =>
                        handleLikeReplies(replyThread)}
                    >
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
                      <span>{replyThread.insights?.likes || 0}</span>
                    </button>
                    <button
                      type="button"
                      class="post-action bookmark"
                      class:active={replyThread.isBookmarked}
                      on:click|stopPropagation={() =>
                        handleBookmarkReplies(replyThread)}
                      aria-label="Bookmark Post"
                    >
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
                          d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v16z"
                        ></path>
                      </svg>
                    </button>
                  </div>
                </div>
              {/each}
            </div>
          </div>
        </div>
      </div>
    </div>
  {/if}
{/if}

<style>
  .thread-detail-container {
    display: flex;
    width: 50%;
    height: 100vh;
  }

  .content-area {
    flex: 1;
    border-right: 1px solid var(--border-color);
    width: 100%;
    height: 100vh;
    margin: 0 auto;
    overflow-y: auto;
  }

  .header {
    position: sticky;
    top: 0;
    background-color: var(--card-background);
    backdrop-filter: blur(12px);
    z-index: 10;
    border-bottom: 1px solid var(--border-color);
    display: flex;
    align-items: center;
    padding: 12px 16px;
    gap: 20px;
  }

  .back-button {
    background: none;
    border: none;
    color: var(--text-primary);
    cursor: pointer;
    width: 34px;
    height: 34px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;
  }

  .back-button:hover {
    background-color: var(--hover-background);
  }

  .header h1 {
    font-size: 20px;
    font-weight: 700;
    color: var(--text-primary);
    margin: 0;
  }

  .thread-post {
    display: flex;
    padding: 16px;
    border-bottom: 1px solid var(--border-color);
  }

  .post-avatar {
    margin-right: 12px;
    flex-shrink: 0;
  }

  .thread-post .post-avatar img {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    object-fit: cover;
  }
  .post-avatar img {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
  }

  .avatar-button {
    cursor: pointer;
    border-radius: 50%;
    transition: opacity 0.2s;
    display: block;
    width: fit-content;
  }

  .avatar-button:hover {
    opacity: 0.8;
  }

  .avatar-button img {
    display: block;
  }

  .media-button {
    cursor: pointer;
    transition: opacity 0.2s;
    display: block;
    width: 100%;
    height: 100%;
  }

  .media-button:hover {
    opacity: 0.9;
  }

  .media-button img,
  .media-button video {
    display: block;
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .post-content {
    flex: 1;
    min-width: 0;
  }

  .post-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 4px;
  }

  .post-user-info {
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .post-reply-info {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 4px;
    min-width: 0;
    line-height: 1.3;
  }

  .thread-post .post-user-info {
    display: flex;
    flex-direction: column;
  }

  .user-name-container {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .post-user-name {
    font-weight: 700;
    color: var(--text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .thread-post .post-user-name {
    font-size: 16px;
  }

  .verified-badge {
    display: flex;
    align-items: center;
    flex-shrink: 0;
  }

  .post-username,
  .post-date {
    color: var(--text-tertiary);
    white-space: nowrap;
    font-size: 14px;
  }

  .thread-post .post-username {
    font-size: 15px;
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
    padding: 4px;
    margin-left: 8px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition:
      background-color 0.2s,
      color 0.2s;
    width: 30px;
    height: 30px;
    flex-shrink: 0;
  }

  .post-more:hover {
    background-color: rgba(var(--button-primary-rgb, 29, 155, 240), 0.1);
    color: var(--button-primary);
  }

  .post-text {
    color: var(--text-primary);
    margin-bottom: 12px;
    word-wrap: break-word;
    line-height: 1.5;
    font-size: 15px;
  }

  .thread-post .post-text {
    font-size: 17px;
  }

  .post-time {
    color: var(--text-tertiary);
    font-size: 15px;
    margin: 16px 0;
    padding-bottom: 16px;
    border-bottom: 1px solid var(--border-color);
  }

  .post-images {
    border-radius: 16px;
    overflow: hidden;
    margin-top: 8px;
    margin-bottom: 12px;
    display: grid;
    gap: 2px;
    border: 1px solid var(--border-color);
  }

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
    position: relative;
    overflow: hidden;
    background-color: var(--hover-background);
  }

  .image-container img {
    display: block;
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: opacity 0.3s ease;
  }

  .post-actions {
    display: flex;
    justify-content: space-between;
    max-width: 425px;
    margin-top: 12px;
  }

  .post-action {
    background: none;
    border: none;
    color: var(--text-tertiary);
    display: flex;
    align-items: center;
    gap: 6px;
    cursor: pointer;
    transition: color 0.2s;
    font-size: 13px;
    padding: 6px;
    border-radius: 4px;
  }

  .post-action.reply:hover {
    color: var(--button-primary);
    background-color: rgba(var(--button-primary-rgb, 29, 155, 240), 0.1);
  }

  .post-action.reply.active {
    color: var(--button-primary);
    background-color: rgba(var(--button-primary-rgb, 29, 155, 240), 0.5);
  }

  .post-action.like:hover {
    color: red;
    background-color: rgba(var(--button-primary-rgb, 200, 15, 24), 0.1);
  }
  .post-action.like.active {
    color: red;
    background-color: rgba(var(--button-primary-rgb, 200, 15, 24), 0.7);
  }

  .post-action.repost:hover {
    color: green;
    background-color: rgba(var(--button-primary-rgb, 29, 240, 40), 0.1);
  }

  .post-action.repost.active {
    color: green;
    background-color: rgba(var(--button-primary-rgb, 29, 240, 40), 0.7);
  }

  .post-action.bookmark:hover {
    color: var(--button-primary);
    background-color: rgba(var(--button-primary-rgb, 29, 155, 240), 0.1);
  }

  .post-action.bookmark.active {
    color: var(--button-primary);
    background-color: rgba(var(--button-primary-rgb, 29, 155, 240), 0.7);
  }

  .post-action svg {
    width: 18px;
    height: 18px;
  }

  .reply-to {
    color: var(--text-tertiary);
    font-size: 14px;
    margin-bottom: 4px;
  }

  .reply-username {
    color: var(--button-primary);
  }

  .no-replies {
    padding: 32px 16px;
    text-align: center;
    color: var(--text-tertiary);
  }

  .posts {
    display: flex;
    flex-direction: column;
  }

  .post {
    display: flex;
    padding: 12px 16px;
    border-bottom: 1px solid var(--border-color);
    transition: background-color 0.2s;
    cursor: pointer;
  }

  .post:hover {
    background-color: var(--hover-background);
  }

  .loading-sentinel {
    min-height: 100px;
    padding: 20px;
    text-align: center;
    color: var(--text-secondary);
  }

  .skeleton-loading {
    width: 90%;
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

  .pulse {
    animation: pulse 1.5s infinite;
  }
  .post-video {
    position: relative;
    width: 100%;
    height: 280px;
    background-color: #000;
    overflow: hidden;
  }
  /* Media Overlay Styles */
  .media-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    z-index: 1000;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: var(--background);
  }

  .media-overlay:focus {
    outline: none;
  }

  .media-overlay-backdrop {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    cursor: pointer;
  }

  .media-overlay-content {
    position: relative;
    width: 100%;
    height: 100%;
    display: flex;
    max-width: 1400px;
    max-height: 100vh;
  }

  .overlay-close-button {
    position: absolute;
    top: 16px;
    right: 16px;
    background: rgba(0, 0, 0, 0.6);
    border: none;
    color: white;
    cursor: pointer;
    width: 44px;
    height: 44px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;
    z-index: 1001;
  }

  .overlay-close-button:hover {
    background: rgba(0, 0, 0, 0.8);
  }

  /* Left side: Media viewer */
  .media-overlay-left {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
  }

  .media-overlay-viewer {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    padding: 60px 20px 20px;
  }

  .media-nav-button {
    position: absolute;
    background: rgba(0, 0, 0, 0.6);
    border: none;
    color: white;
    cursor: pointer;
    width: 44px;
    height: 44px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition:
      background-color 0.2s,
      opacity 0.2s;
    z-index: 1001;
  }

  .media-nav-button:hover:not(.disabled) {
    background: rgba(0, 0, 0, 0.8);
  }

  .media-nav-button.disabled {
    opacity: 0.3;
    cursor: not-allowed;
  }

  .media-nav-button.prev {
    left: 20px;
  }

  .media-nav-button.next {
    right: 20px;
  }

  .media-overlay-media {
    max-width: 100%;
    max-height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .overlay-media-item {
    max-width: 100%;
    max-height: 100%;
    object-fit: contain;
    border-radius: 8px;
  }

  .overlay-video {
    width: 100%;
    height: auto;
  }

  .media-counter {
    position: absolute;
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
    background: rgba(0, 0, 0, 0.6);
    color: white;
    padding: 8px 16px;
    border-radius: 20px;
    font-size: 14px;
    font-weight: 500;
  }

  /* Right side: Thread interactions */
  .media-overlay-right {
    width: 400px;
    background: var(--card-background);
    border-left: 1px solid var(--border-color);
    display: flex;
    flex-direction: column;
    height: 100vh;
  }

  .overlay-sidebar {
    flex: 1;
    overflow-y: auto;
    padding: 0;
  }

  .overlay-original-thread {
    padding: 16px;
    border-bottom: 1px solid var(--border-color);
  }

  .overlay-thread-header {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    margin-bottom: 12px;
  }

  .overlay-thread-header .post-avatar img {
    width: 40px;
    height: 40px;
  }
  .post-avatar img:hover {
    opacity: 0.8;
  }

  .overlay-user-info {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 4px;
    min-width: 0;
    line-height: 1.3;
  }

  .overlay-post-text {
    color: var(--text-primary);
    margin-bottom: 12px;
    word-wrap: break-word;
    line-height: 1.5;
    font-size: 15px;
  }

  .overlay-post-actions {
    display: flex;
    justify-content: space-between;
    max-width: 320px;
    margin-top: 8px;
  }
  /* Reply composition */
  .overlay-reply-compose {
    padding: 16px;
    border-bottom: 1px solid var(--border-color);
    background: var(--card-background);
  }

  .overlay-compose-reply {
    padding: 16px;
    border-bottom: 1px solid var(--border-color);
    background: var(--card-background);
  }

  .overlay-compose-reply h3 {
    font-size: 16px;
    font-weight: 700;
    color: var(--text-primary);
    margin: 0 0 16px 0;
    padding-bottom: 12px;
    border-bottom: 1px solid var(--border-color);
  }

  .reply-compose-header {
    display: flex;
    gap: 12px;
    align-items: flex-start;
  }

  .reply-textarea {
    flex: 1;
    border: none;
    background: transparent;
    color: var(--text-primary);
    font-size: 15px;
    line-height: 1.5;
    resize: none;
    min-height: 60px;
    font-family: inherit;
    outline: none;
  }

  .reply-textarea::placeholder {
    color: var(--text-muted);
  }

  .reply-compose-actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 12px;
  }

  .reply-submit-btn {
    background: var(--accent-color);
    color: white;
    border: none;
    border-radius: 20px;
    padding: 8px 20px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .reply-submit-btn:hover:not(:disabled) {
    background: var(--accent-hover);
  }

  .reply-submit-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  /* Replies list */
  .overlay-replies {
    flex: 1;
  }

  .overlay-reply-item {
    padding: 16px;
    border-bottom: 1px solid var(--border-color);
  }

  .reply-header {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    margin-bottom: 8px;
  }

  .reply-user-info {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 4px;
    min-width: 0;
    line-height: 1.3;
  }

  .reply-text {
    color: var(--text-primary);
    margin-bottom: 12px;
    word-wrap: break-word;
    line-height: 1.5;
    font-size: 15px;
  }

  .reply-images {
    display: grid;
    gap: 8px;
    margin-bottom: 12px;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }

  .reply-image {
    width: 100%;
    height: 150px;
    object-fit: cover;
    border-radius: 8px;
    cursor: pointer;
    transition: opacity 0.2s;
  }
  .reply-image:hover {
    opacity: 0.9;
  }

  .reply-image-button {
    cursor: pointer;
    transition: opacity 0.2s;
    display: block;
    width: 100%;
    height: 150px;
    border-radius: 8px;
    overflow: hidden;
  }

  .reply-image-button:hover {
    opacity: 0.9;
  }

  .reply-image-button .reply-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 0;
  }

  .reply-actions {
    display: flex;
    justify-content: space-between;
    max-width: 300px;
  }

  /* Responsive adjustments */
  @media (max-width: 1000px) {
    .media-overlay-content {
      flex-direction: column;
    }

    .media-overlay-left {
      flex: 1;
      min-height: 60vh;
    }

    .media-overlay-right {
      width: 100%;
      height: 40vh;
      border-left: none;
      border-top: 1px solid var(--border-color);
    }
  }

  /* Make images in post clickable */
  .image-container img,
  .post-video {
    cursor: pointer;
    transition: opacity 0.2s;
  }

  .image-container img:hover,
  .post-video:hover {
    opacity: 0.9;
  }

  @keyframes pulse {
    0% {
      opacity: 0.6;
    }
    50% {
      opacity: 1;
    }
    100% {
      opacity: 0.6;
    }
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }

  @media (max-width: 768px) {
    .content-area {
      border-right: none;
      width: 100%;
    }
    .post-actions {
      max-width: 100%;
    }
  }

  @media (max-width: 500px) {
    .post-user-info {
      font-size: 14px;
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
    .thread-post .post-text {
      font-size: 16px;
    }
  }
  @media (max-width: 1024px) {
    .thread-detail-container {
      flex-direction: column;
    }

    .content-area {
      border-right: none;
      height: 100vh;
    }
  }

  @media (max-width: 768px) {
    .content-area {
      height: 100vh;
    }
  }
</style>
