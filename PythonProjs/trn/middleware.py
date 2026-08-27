from fastapi.middleware.cors import CORSMiddleware

def setup_middleware(app):
    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],  # Allow all origins (modify for security)
        allow_credentials=True,
        allow_methods=["*"],  # Allow all HTTP methods
        allow_headers=["*"]   # Allow all headers
    )
