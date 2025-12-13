import os
from transformers import AutoConfig, AutoTokenizer

# --- CONFIGURATION ---
# PASTE YOUR FOLDER PATH HERE (Where model.safetensors is located)
FOLDER_PATH = r"C:/Users/DELL/OneDrive/Desktop/GenAI Project/final/codet5-finetuned-corrected"

# The base model to borrow config/tokenizer from
BASE_MODEL = "Salesforce/codet5-small"

def repair_folder():
    print(f"🔧 Repairing model folder: {FOLDER_PATH}")
    
    if not os.path.exists(FOLDER_PATH):
        print("❌ Error: Folder not found!")
        return

    # 1. Download and Save Config (The Skeleton)
    print("⬇️ Downloading missing config.json...")
    try:
        config = AutoConfig.from_pretrained(BASE_MODEL)
        config.save_pretrained(FOLDER_PATH)
        print("✅ config.json saved.")
    except Exception as e:
        print(f"❌ Failed to save config: {e}")

    # 2. Download and Save Tokenizer (The Dictionary)
    print("⬇️ Downloading missing tokenizer files...")
    try:
        tokenizer = AutoTokenizer.from_pretrained(BASE_MODEL)
        tokenizer.save_pretrained(FOLDER_PATH)
        print("✅ Tokenizer files saved.")
    except Exception as e:
        print(f"❌ Failed to save tokenizer: {e}")

    # 3. Verify
    print("\n--- Verification ---")
    files = os.listdir(FOLDER_PATH)
    required = ["config.json", "model.safetensors"]
    
    missing = [f for f in required if f not in files]
    
    if not missing:
        print("🚀 SUCCESS! Your model folder is now valid.")
        print("You can now run 'streamlit run app.py'")
    else:
        print(f"⚠️ Still missing: {missing}")
        print("Ensure 'model.safetensors' is actually in this folder!")

if __name__ == "__main__":
    repair_folder()