
def generate_data_sql(num_owners=1000, pets_per_owner=10):
    sql_statements = []
    
    # Keep existing data
    sql_statements.append("""
INSERT IGNORE INTO vets VALUES (2, 'Helen', 'Leary');
INSERT IGNORE INTO vets VALUES (3, 'Linda', 'Douglas');
INSERT IGNORE INTO vets VALUES (4, 'Rafael', 'Ortega');
INSERT IGNORE INTO vets VALUES (5, 'Henry', 'Stevens');
INSERT IGNORE INTO vets VALUES (6, 'Sharon', 'Jenkins');

INSERT IGNORE INTO specialties VALUES (1, 'radiology');
INSERT IGNORE INTO specialties VALUES (2, 'surgery');
INSERT IGNORE INTO specialties VALUES (3, 'dentistry');

INSERT IGNORE INTO vet_specialties VALUES (2, 1);
INSERT IGNORE INTO vet_specialties VALUES (3, 2);
INSERT IGNORE INTO vet_specialties VALUES (3, 3);
INSERT IGNORE INTO vet_specialties VALUES (4, 2);
INSERT IGNORE INTO vet_specialties VALUES (5, 1);

INSERT IGNORE INTO types VALUES (1, 'cat');
INSERT IGNORE INTO types VALUES (2, 'dog');
INSERT IGNORE INTO types VALUES (3, 'lizard');
INSERT IGNORE INTO types VALUES (4, 'snake');
INSERT IGNORE INTO types VALUES (5, 'bird');
INSERT IGNORE INTO types VALUES (6, 'hamster');

INSERT IGNORE INTO owners VALUES (1, 'George', 'Franklin', '110 W. Liberty St.', 'Madison', '6085551023');
INSERT IGNORE INTO owners VALUES (2, 'Betty', 'Davis', '638 Cardinal Ave.', 'Sun Prairie', '6085551749');
INSERT IGNORE INTO owners VALUES (3, 'Eduardo', 'Rodriquez', '2693 Commerce St.', 'McFarland', '6085558763');
INSERT IGNORE INTO owners VALUES (4, 'Harold', 'Davis', '563 Friendly St.', 'Windsor', '6085553198');
INSERT IGNORE INTO owners VALUES (5, 'Peter', 'McTavish', '2387 S. Fair Way', 'Madison', '6085552765');
INSERT IGNORE INTO owners VALUES (6, 'Jean', 'Coleman', '105 N. Lake St.', 'Monona', '6085552654');
INSERT IGNORE INTO owners VALUES (7, 'Jeff', 'Black', '1450 Oak Blvd.', 'Monona', '6085555387');
INSERT IGNORE INTO owners VALUES (8, 'Maria', 'Escobito', '345 Maple St.', 'Madison', '6085557683');
INSERT IGNORE INTO owners VALUES (9, 'David', 'Schroeder', '2749 Blackhawk Trail', 'Madison', '6085559435');
INSERT IGNORE INTO owners VALUES (10, 'Carlos', 'Estaban', '2335 Independence La.', 'Waunakee', '6085555487');
    """)
    
    # Generate Owners
    sql_statements.append("-- Insert Owners")
    for i in range(11, num_owners + 11):
        sql_statements.append(
            f"INSERT IGNORE INTO owners (id, first_name, last_name, address, city, telephone) VALUES ({i}, 'Owner{i}', 'LastName{i}', 'Address{i}', 'City{i}', '0100000{i:04d}');"
        )
    
    # Generate Pets
    sql_statements.append("\n-- Insert Pets")
    pet_id = 14  # Continue from existing pets
    for owner_id in range(11, num_owners + 11):
        for j in range(1, pets_per_owner + 1):
            pet_name = f"Pet{j}_Owner{owner_id}"
            birth_date = f"2020-{(j % 12) + 1:02d}-{(j % 28) + 1:02d}"  # Distributes birth dates
            type_id = (j % 6) + 1  # Randomized between 1 and 6
            sql_statements.append(
                f"INSERT IGNORE INTO pets (id, name, birth_date, type_id, owner_id) VALUES ({pet_id}, '{pet_name}', '{birth_date}', {type_id}, {owner_id});"
            )
            pet_id += 1
    
    # Keep existing visits and users
    sql_statements.append("""
INSERT IGNORE INTO visits VALUES (1, 7, '2010-03-04', 'rabies shot');
INSERT IGNORE INTO visits VALUES (2, 8, '2011-03-04', 'rabies shot');
INSERT IGNORE INTO visits VALUES (3, 8, '2009-06-04', 'neutered');
INSERT IGNORE INTO visits VALUES (4, 7, '2008-09-04', 'spayed');

INSERT IGNORE INTO users(username,password,enabled) VALUES ('admin','{noop}admin', true);

INSERT IGNORE INTO roles (username, role) VALUES ('admin', 'ROLE_OWNER_ADMIN');
INSERT IGNORE INTO roles (username, role) VALUES ('admin', 'ROLE_VET_ADMIN');
INSERT IGNORE INTO roles (username, role) VALUES ('admin', 'ROLE_ADMIN');
    """)
    
    return "\n".join(sql_statements)

# Save to a file
data_sql_content = generate_data_sql()
with open("./src/main/resources/db/mysql/data.sql", "w") as f:
    f.write(data_sql_content)

print("data.sql file generated successfully!")

