import streamlit as st
import google.generativeai as genai
import time

# --- CONFIGURATION ---
# 🔴 PASTE YOUR API KEY HERE
API_KEY = "AIzaSyCRjCYhFpeb55hnbbRRCOr_5rEvc27IoIU"

# Configure the API
genai.configure(api_key=API_KEY)

st.set_page_config(page_title="CodeScribe AI", layout="wide")

# --- UI LAYOUT ---
st.title("🤖 CodeScribe: AI Code Commenter")
st.markdown("Enter a code snippet below, and the AI will generate inline comments and a summary.")

# Input Area
instruction = st.text_input("Instruction", value="Add inline comments and a summary to this code.")
code_input = st.text_area("Source Code", height=200, placeholder="def my_function(x): ...")

def get_working_model():
    """
    Dynamically finds ANY model that works for this API key.
    """
    try:
        # Search through all available models
        for m in genai.list_models():
            if 'generateContent' in m.supported_generation_methods:
                print(f"✅ Found working model: {m.name}")
                return genai.GenerativeModel(m.name)
    except Exception as e:
        st.error(f"Error listing models: {e}")
        return None
    
    return None

def get_gemini_response(instruction, code):
    try:
        # Auto-select the best model
        model = get_working_model()
        
        if not model:
            return "Error: No compatible models found for your API Key. Check your Google AI Studio permissions."

        prompt = f"""
        You are an AI code assistant.
        Task: {instruction}
        
        Input Code:
        {code}
        
        Please provide the python code with inline comments explaining each step, 
        followed by a brief summary at the end. 
        Do not use markdown formatting like ```python. Just give the raw code text.
        """
        
        response = model.generate_content(prompt)
        return response.text
    except Exception as e:
        return f"Error: {str(e)}"

# --- MAIN LOGIC ---
if st.button("Generate Comments"):
    if not code_input:
        st.warning("⚠️ Please enter some code first.")
    else:
        # Fake progress bar for demo effect
        progress_text = "Loading model components..."
        my_bar = st.progress(0, text=progress_text)

        for percent_complete in range(100):
            time.sleep(0.01) 
            my_bar.progress(percent_complete + 1, text="Analyzing syntax trees...")
            
        my_bar.empty()

        with st.spinner("Generating explanations..."):
            result = get_gemini_response(instruction, code_input)
            
            # Display Result
            st.subheader("✅ Commented Code:")
            st.code(result, language='python')
            
            st.success("Processing complete.")