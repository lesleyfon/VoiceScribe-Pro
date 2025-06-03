1. Activate the virtual environment
```bash
source venv/bin/activate
```
2. `cd` into the `ml` directory
3. Run the app
```bash
 uvicorn main:app --host 0.0.0.0 --port 9090
```
  * To run the app in watch mode, use `uvicorn main:app --host 0.0.0.0 --port 9090 --reload`
