from typing import Optional
import requests


def get_latest_otp(base_url: str, email: str, timeout: int = 10) -> Optional[str]:
    """Fetch the latest OTP for an email from the test hook endpoint.

    Returns the OTP string or None if not found/endpoint disabled.
    """
    try:
        url = f"{base_url.rstrip('/')}/auth/test/get-latest-otp"
        resp = requests.get(url, params={"email": email}, timeout=timeout)
        if resp.status_code != 200:
            return None
        data = resp.json() or {}
        return data.get("otp")
    except Exception:
        return None


