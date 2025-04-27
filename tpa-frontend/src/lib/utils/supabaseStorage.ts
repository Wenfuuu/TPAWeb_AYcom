// import { createClient } from "@supabase/supabase-js";

// Initialize Supabase client
// const supabaseUrl = "https://nyijguoyjemxmdyzzdry.supabase.co";
// const supabaseKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im55aWpndW95amVteG1keXp6ZHJ5Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NDU2NTI4MTMsImV4cCI6MjA2MTIyODgxM30.cGZzeD6ojLZjxPu9SpFGkjnkPCAnUNQy8GFCXkLD8uk";
// const supabase = createClient(supabaseUrl, supabaseKey);

export async function uploadRegistrationFile(
  file: File,
  fileType: "profile" | "banner"
): Promise<string | null> {

  const formData = new FormData();
  formData.append("file", file);

  const response = await fetch(`http://localhost:8080/api/users/upload/${fileType}`, {
    method: "POST",
    body: formData,
  });

  const data = await response.json();
  return data.url;
}

export async function uploadProfilePicture(file: File): Promise<string | null> {
  return uploadRegistrationFile(file, "profile");
}

export async function uploadBannerImage(file: File): Promise<string | null> {
  return uploadRegistrationFile(file, "banner");
}