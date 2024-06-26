import streamlit as st
import pandas as pd
import requests
import os


# Define a function to fetch data from the API
def fetch_data(limit, page, isEncrypted, groupDuplicates):
    API_SERVER_ENDPOINT = os.getenv("API_SERVER_ENDPOINT", "http://localhost:8080/login-data")
    # Replace 'YOUR_API_ENDPOINT' with your actual API endpoint
    response = requests.get('{0}?limit={1}&page={2}&isEncrypted={3}&groupDuplicates={4}'.format(API_SERVER_ENDPOINT, limit, page, isEncrypted, groupDuplicates))
    if response.status_code == 200:
        return response.json()
    else:
        st.error(f"Error fetching data: {response.status_code}")
        return []


# Define a function to load data and manage masking
def load_data(limit, page, isEncrypted, groupDuplicates):
    data = fetch_data(limit, page, isEncrypted, groupDuplicates)
    df = pd.DataFrame(data)
    return df


def init_session():
    if 'limit' not in st.session_state:
        st.session_state['limit'] = 5

    if 'page' not in st.session_state:
        st.session_state['page'] = 1

    if 'isEncrypted' not in st.session_state:
        st.session_state['isEncrypted'] = True

    if 'groupDuplicates' not in st.session_state:
        st.session_state['groupDuplicates'] = False



# Main Streamlit app function
def main():
    st.title('User Logins')

    col1, col2, col3, col4 = st.columns([1, 1, 1, 1])
    with col1:
        limit = st.slider("Select records limit", 0, 100,  st.session_state['limit'])
        st.session_state['limit'] = limit
        st.write("Records limit: ",  st.session_state['limit'])

    with col2:
        page = st.number_input("Select the page", value=st.session_state['page'])
        st.session_state['page'] = page
        st.write("Page ",  st.session_state['page'])

    with col3:
        st.write("\n\n\n")
        if st.button("Show Encrypted", type="primary"):
            if st.session_state['isEncrypted']:
                st.session_state['isEncrypted'] = False
            else:
                st.session_state['isEncrypted'] = True
        st.write("Data Encrypted: ", st.session_state['isEncrypted'])

    with col4:
        st.write("\n\n\n")
        if st.button("Group Duplicates", type="primary"):
            if st.session_state['groupDuplicates']:
                st.session_state['groupDuplicates'] = False
            else:
                st.session_state['groupDuplicates'] = True

        st.write("Group Duplicates: ", st.session_state['groupDuplicates'])

    if st.button("Fetch Data"):
        # Load data
        data_load_state = st.text('Loading data...')
        df = load_data(st.session_state['limit'], st.session_state['page'], st.session_state['isEncrypted'], st.session_state['groupDuplicates'])
        data_load_state.text('')

        st.dataframe(data=df)


init_session()
# Refresh data every 1 minute
main()

