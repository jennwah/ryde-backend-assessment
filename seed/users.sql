INSERT INTO users.users (id, name, date_of_birth, address, description, location)
VALUES ('4f4ba666-37a7-4873-baac-d4404aa3a03d', 'Jenn Wah', '1997-12-17', 'Kuala Lumpur', 'just a chilled guy', public.ST_SetSRID(public.ST_MakePoint(100.22, 45.22), 4326)),
       ('dbf88f17-fe5b-4027-8831-c1acbedd0770', 'John Doe', '1988-12-11', 'Penang', 'second user', public.ST_SetSRID(public.ST_MakePoint(95.22, 44.22), 4326)),
       ('7195dcd0-9e78-4eb9-9d63-5af46c4739da','Sabrina', '1947-01-17', 'Singapore', 'third user', public.ST_SetSRID(public.ST_MakePoint(105.22, 48.22), 4326));

INSERT INTO users.friends (user_id, friend_id)
VALUES ('4f4ba666-37a7-4873-baac-d4404aa3a03d', '7195dcd0-9e78-4eb9-9d63-5af46c4739da'),
       ('7195dcd0-9e78-4eb9-9d63-5af46c4739da', 'dbf88f17-fe5b-4027-8831-c1acbedd0770');