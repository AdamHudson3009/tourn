interface CookieOptions {
  days?: number;                       // expiry in days (omit for a session cookie)
  path?: string;                       // default "/"
  domain?: string;
  secure?: boolean;
  sameSite?: "Strict" | "Lax" | "None";
}

export function setCookie(name: string, value: string, options: CookieOptions = {}): void {
  const { days, path = "/", domain, secure, sameSite = "Lax" } = options;

  let cookie = `${encodeURIComponent(name)}=${encodeURIComponent(value)}`;

  if (days !== undefined) {
    const expires = new Date(Date.now() + days * 24 * 60 * 60 * 1000);
    cookie += `; Expires=${expires.toUTCString()}`;
  }

  cookie += `; Path=${path}`;
  if (domain) cookie += `; Domain=${domain}`;
  if (secure || sameSite === "None") cookie += "; Secure"; // SameSite=None requires Secure
  cookie += `; SameSite=${sameSite}`;

  document.cookie = cookie;
}

export function getCookie(name: string): string | null {
  const target = `${encodeURIComponent(name)}=`;

  for (const part of document.cookie.split("; ")) {
    if (part.startsWith(target)) {
      return decodeURIComponent(part.slice(target.length));
    }
  }
  return null;
}

export function clearCookie(
  name: string,
  options: Pick<CookieOptions, "path" | "domain"> = {}
): void {
  const { path = "/", domain } = options;

  let cookie = `${encodeURIComponent(name)}=; Expires=Thu, 01 Jan 1970 00:00:00 GMT; Max-Age=0; Path=${path}`;
  if (domain) cookie += `; Domain=${domain}`;

  document.cookie = cookie;
}