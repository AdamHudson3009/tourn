from fastapi import FastAPI
import logging
from middleware import setup_middleware  # Import middleware setup
from leagues import leagueRtr  # Import the API router
from tcd import tcdRtr  # Import the API router
from leagues import leagueRtr  # Import the API router
from plyr import plyrRtr  # Import the API router
from tbrk import tbrkRtr  # Import the API router
from trnTable import trnTableRtr
from cmm import cmmRtr
from grammar import grammarRtr 
from lllsch import lllschRtr
from lllextra import lllExtraRtr

app = FastAPI()
logging.basicConfig(level=logging.INFO)

# Apply middleware
setup_middleware(app)

# Include routes
app.include_router(leagueRtr)
app.include_router(tcdRtr)
app.include_router(plyrRtr)
app.include_router(tbrkRtr)
app.include_router(trnTableRtr)
app.include_router(cmmRtr)
app.include_router(grammarRtr)
app.include_router(lllschRtr)
app.include_router(lllExtraRtr)