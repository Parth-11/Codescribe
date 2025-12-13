# app.py
import streamlit as st
from transformers import AutoTokenizer, AutoModelForSeq2SeqLM
import torch

# --- CONFIGURATION ---
# Point this to your unzipped model folder
MODEL_PATH = r"C:/Users/DELL/OneDrive/Desktop/GenAI Project/final/codet5-finetuned-corrected"


st.set_page_config(page_title="CodeScribe AI", layout="wide")

@st.cache_resource
def load_model():
    try:
        tokenizer = AutoTokenizer.from_pretrained(MODEL_PATH, local_files_only=True)
        model = AutoModelForSeq2SeqLM.from_pretrained(MODEL_PATH, local_files_only=True)
        return tokenizer, model
    except Exception as e:
        st.error(f"Error loading model: {e}")
        return None, None

st.title("🤖 CodeScribe: AI Code Commenter")
st.markdown("Enter a code snippet below, and the AI will generate inline comments and a summary.")

# Load Model
tokenizer, model = load_model()

# Input Area
instruction = st.text_input("Instruction", value="Add inline comments and a summary to this code.")
code_input = st.text_area("Source Code", height=200, placeholder="def my_function(x): ...")

if st.button("Generate Comments"):
    if not code_input:
        st.warning("Please enter some code first.")
    elif tokenizer and model:
        with st.spinner("Generating explanations..."):
            # --- CRITICAL: MATCHING THE TRAINING FORMAT ---
            # We must format the input exactly as we did in the Colab training script
            input_text = f"Instruction: {instruction}\nInput: {code_input}"
            
            # Tokenize
            inputs = tokenizer(input_text, return_tensors="pt", max_length=256, truncation=True)
            
            # Generate
            summary_ids = model.generate(
                inputs["input_ids"],
                max_length=256, 
                num_beams=4,
                early_stopping=True
            )
            
            # Decode
            result = tokenizer.decode(summary_ids[0], skip_special_tokens=True)
            
            # Display Result
            st.subheader("✅ Commented Code:")
            st.code(result, language='python')