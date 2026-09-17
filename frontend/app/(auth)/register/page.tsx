import type { Metadata } from "next";
import { RegisterForm } from "@/components/auth/register-form";

export const metadata: Metadata = {
  title: "Create account — CampusX",
  description: "Join CampusX and discover college events near you",
};

export default function RegisterPage() {
  return <RegisterForm />;
}