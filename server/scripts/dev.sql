INSERT INTO categories(tag) VALUES('Bride & Groom');
INSERT INTO categories(tag) VALUES('Groom''s Relatives');

INSERT INTO invitations(category_id,private_id,greeting,maximum_guest_count,status,notes,mobile_phone_number)
VALUES(1,'12345-67890','Mitten Lin',2,'NS','','94444444');
INSERT INTO invitations(category_id,private_id,greeting,maximum_guest_count,status,notes,mobile_phone_number)
VALUES(2,'54321-09876','Mr Tan Kok Leong & Family',4,'ST','Need 1 baby chair','91501234');
INSERT INTO invitations(category_id,private_id,greeting,maximum_guest_count,status,notes,mobile_phone_number)
VALUES(2,'32145-32145','Mr Steven Tan & Family',3,'ST','','94412321');

INSERT INTO rsvps(invitation_private_id,full_name,attending,guest_count,special_diet,remarks,mobile_phone_number)
VALUES('54321-09876','Mr Tan Kok Leong & Family',true,3,false,'Happy wedding!','91501234');

