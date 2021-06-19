#%%
import pandas as pd

# GTFSを読み込みマージ
trips_df = pd.read_csv('./gtfs/trips.txt')
stops_df = pd.read_csv('./gtfs/stops.txt')
stop_times_df = pd.read_csv('./gtfs/stop_times.txt')
routes_df = pd.read_csv('./gtfs/routes.txt')
df = pd.merge(trips_df,routes_df)
df = pd.merge(df,stop_times_df)
df = pd.merge(df,stops_df)
df.to_csv('df.csv')
df
# %%
# # %%
# df = df[df['trip_id'] == '10100_01202104011101A00101']
# df['lat'] = df['stop_lat']
# df['lon'] = df['stop_lon']
# df = df.sort_values('departure_time')
# df.to_csv('trip1.csv')
# df
# %%
df = pd.read_parquet('gtfsrt.parquet')
df = df[df['trip_id'] == '10100_01202104011101A00101']
df = df.sort_values('timestamp')
df.to_csv('trip2.csv')
df
# %%

# %%
