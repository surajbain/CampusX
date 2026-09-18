import { apiFetch } from "./client";

export interface Payment {
  id: string;
  registration_id: string;
  user_id: string;
  event_id: string;
  provider: string;
  gateway_payment_id?: string;
  gateway_order_id?: string;
  amount_paise: number;
  currency: string;
  status: "CREATED" | "AUTHORIZED" | "CAPTURED" | "FAILED" | "REFUNDED";
  verified_at?: string;
  created_at: string;
  updated_at: string;
}

export interface CreatePaymentResponse {
  payment_id: string;
  gateway_order_id: string;
  amount_paise: number;
  currency: string;
  provider: string;
  public_key?: string;
  idempotency_key: string;
}

export interface VerifyPaymentInput {
  payment_id: string;
  gateway_order_id: string;
  gateway_payment_id: string;
  signature: string;
}

// apiFetch returns {data, meta} wrapper — unwrap
function unwrap<T>(res: unknown): T {
  return (res as { data: T }).data;
}

export async function createPayment(
  registrationId: string
): Promise<CreatePaymentResponse> {
  const res = await apiFetch<CreatePaymentResponse>("/api/v1/payments/create", {
    method: "POST",
    body: JSON.stringify({ registration_id: registrationId }),
  });
  return unwrap<CreatePaymentResponse>(res);
}

export async function verifyPayment(
  input: VerifyPaymentInput
): Promise<Payment> {
  const res = await apiFetch<Payment>("/api/v1/payments/verify", {
    method: "POST",
    body: JSON.stringify(input),
  });
  return unwrap<Payment>(res);
}

export async function getPayment(id: string): Promise<Payment> {
  const res = await apiFetch<Payment>(`/api/v1/payments/${id}`);
  return unwrap<Payment>(res);
}