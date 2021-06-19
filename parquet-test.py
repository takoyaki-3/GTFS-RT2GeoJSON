#%%
import pandas as pd

df = pd.read_parquet('gtfsrt.parquet')
df
# %%
df.to_csv('vps.csv')
# %%
