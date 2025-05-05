// import { createClient } from "@supabase/supabase-js";

// Initialize Supabase client
// const supabaseUrl = "https://nyijguoyjemxmdyzzdry.supabase.co";
// const supabaseKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im55aWpndW95amVteG1keXp6ZHJ5Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDU2NTI4MTMsImV4cCI6MjA2MTIyODgxM30.cGZzeD6ojLZjxPu9SpFGkjnkPCAnUNQy8GFCXkLD8uk";
// const supabase = createClient(supabaseUrl, supabaseKey);

export async function uploadMedia(
  file: File,
  fileType: "image" | "video"
): Promise<string | null> {
  const formData = new FormData();
  formData.append("file", file);

  const endpoint = fileType === "image" ? "upload/image" : "upload-video/video";

  const response = await fetch(`http://localhost:8080/api/users/${endpoint}`, {
    method: "POST",
    body: formData,
  });

  if (!response.ok) {
    console.error("Upload failed:", response.statusText);
    return null;
  }

  const data = await response.json();
  return data.url;
}

export async function uploadImage(file: File): Promise<string | null> {
  return uploadMedia(file, "image");
}

export async function uploadVideo(file: File): Promise<string | null> {
  return uploadMedia(file, "video");
}

export async function uploadRegistrationImage(
  file: File,
  fileType: "profile" | "banner"
): Promise<string | null> {
  const formData = new FormData();
  formData.append("file", file);

  const response = await fetch(
    `http://localhost:8080/api/users/upload/${fileType}`,
    {
      method: "POST",
      body: formData,
    }
  );

  const data = await response.json();
  return data.url;
}

export async function uploadProfilePicture(file: File): Promise<string | null> {
  return uploadRegistrationImage(file, "profile");
}

export async function uploadBannerImage(file: File): Promise<string | null> {
  return uploadRegistrationImage(file, "banner");
}
